package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Longreader/go-shortener-url.git/internal/repository"
	"github.com/Longreader/go-shortener-url.git/internal/tools"

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

	return st, nil
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
			if err == sql.ErrNoRows {
				continue
			} else {
				return "", repository.ErrURLAlreadyExists
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
		`SELECT * FROM links WHERE id=$1`,
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

func (db *PsqlStorage) Ping(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()
	err := db.db.PingContext(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}
