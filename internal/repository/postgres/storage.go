package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Longreader/go-shortener-url.git/internal/repository"
	"github.com/Longreader/go-shortener-url.git/internal/tools"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PsqlStorage struct {
	db *sqlx.DB
}

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
) (url repository.URL, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	row := st.db.QueryRowContext(
		ctx,
		`SELECT url FROM links WHERE id=$1`,
		id,
	)

	err = row.Scan(&url)

	if err == sql.ErrNoRows {
		return "", repository.ErrURLNotFound
	} else if err != nil {
		return "", err
	}
	return url, nil

}

func (st *PsqlStorage) GetAll(
	ctx context.Context,
	user repository.User,
) (data []repository.LinkData, err error) {

	data = make([]repository.LinkData, 0)

	rows, err := st.db.QueryContext(
		ctx,
		`SELECT url, id, user_id FROM links WHERE user_id=$1`,
		user,
	)

	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
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
