package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Longreader/go-shortener-url.git/internal/repository"
	"github.com/Longreader/go-shortener-url.git/internal/tools"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PsqlStorage struct {
	db             *sqlx.DB
	deleteCh       chan repository.LinkData
	deleteWg       sync.WaitGroup
	deleteShutdown chan struct{}
}

const (
	delBufferSize    = 100
	delBufferTimeout = time.Second
	shutdownTimeout  = 15 * time.Second
)

func NewPsqlStorage(dsn string) (*PsqlStorage, error) {

	var err error

	st := &PsqlStorage{}

	st.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = st.db.Ping()
	if err != nil {
		return nil, err
	}

	st.Setup()

	return st, nil
}

func (st *PsqlStorage) Setup() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	st.db.MustExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS links (
			id      varchar(255) NOT NULL UNIQUE,
			url     varchar(255) NOT NULL UNIQUE,
			user_id uuid         NOT NULL
		);`,
	)

}

// Set method for PsqlStorage storage
func (st *PsqlStorage) Set(
	ctx context.Context,
	url repository.URL,
	user repository.User,
) (id repository.ID, err error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	for {
		id, err = tools.RandStringBytes(5)
		if err != nil {
			return "", err
		}
		_, err := st.db.ExecContext(
			ctx,
			`INSERT INTO links (id, url, user_id) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`,
			id, url, user,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else {
				row := st.db.QueryRowContext(
					ctx,
					`SELECT id FROM links WHERE url=$1`,
					url,
				)
				err := row.Scan(&id)
				if err != nil {
					logrus.Printf("Error scan value %s", err)
					return "", err
				}
				return id, repository.ErrURLAlreadyExists
			}
		} else {
			break
		}
	}
	return id, nil
}

// Get method for PsqlStorage storage
func (st *PsqlStorage) Get(
	ctx context.Context,
	id repository.ID,
) (url repository.URL, deleted bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	row := st.db.QueryRowContext(
		ctx,
		`SELECT url FROM links WHERE id=$1`,
		id,
	)

	err = row.Scan(&url, &deleted)

	if err == sql.ErrNoRows {
		return "", false, repository.ErrURLNotFound
	} else if err != nil {
		return "", false, err
	}
	return url, deleted, nil
}

func (st *PsqlStorage) Delete(
	ctx context.Context,
	ids []repository.ID,
	user repository.User,
) error {
	for _, id := range ids {
		st.deleteCh <- repository.LinkData{ID: id, User: user}
	}
	return nil
}

func (st *PsqlStorage) DeleteLink(ctx context.Context, bufferSize int, bufferTimeout time.Duration) {
	ids := make([]repository.ID, 0, bufferSize)
	users := make([]repository.User, 0, bufferSize)

worker:
	for {
		select {

		case <-st.deleteShutdown:
			break worker

		case <-ctx.Done():
			break worker

		default:

		}

		ids = ids[:0]
		users = users[:0]

		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, bufferTimeout)

	loop:
		for {
			select {
			case v := <-st.deleteCh:
				ids = append(ids, v.ID)
				users = append(users, v.User)
				if len(ids) == bufferSize {
					timeoutCancel()
					break loop
				}
			case <-timeoutCtx.Done():
				timeoutCancel()
				break loop
			case <-st.deleteShutdown:
				timeoutCancel()
				break loop
			}
		}

		if len(ids) == 0 {
			continue
		}

		ctxLocal, cancelLocal := context.WithTimeout(ctx, time.Second*10)

		_, err := st.db.ExecContext(
			ctxLocal,
			`UPDATE links SET deleted = TRUE 
             FROM (SELECT UNNEST($1::text[]) AS id, UNNEST($2::uuid[]) AS user) AS data_table
             WHERE links.id = data_table.id AND user_id = data_table.user`,
			ids, users,
		)
		if err != nil {
			log.Printf("update failed: %v", err)
		}

		cancelLocal()
	}

	st.deleteWg.Done()
}

func (st *PsqlStorage) GetAllByUser(
	ctx context.Context,
	user repository.User,
) (data []repository.LinkData, err error) {

	data = make([]repository.LinkData, 0)

	rows, err := st.db.QueryContext(
		ctx,
		`SELECT url, id, user_id FROM links WHERE user_id=$1 and deleted=FALSE`,
		user,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var ld repository.LinkData

		err := rows.Scan(&ld.URL, &ld.ID, &ld.User)
		if err != nil {
			return data, err
		}

		data = append(data, ld)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (st *PsqlStorage) Ping(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()
	err := st.db.PingContext(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (st *PsqlStorage) Close(_ context.Context) error {
	close(st.deleteShutdown)

	c := make(chan struct{})
	go func() {
		defer close(c)
		st.deleteWg.Wait()
	}()
	select {
	case <-c:
	case <-time.After(shutdownTimeout):
		log.Print("storage close timeout exceed")
	}

	return st.db.Close()
}
