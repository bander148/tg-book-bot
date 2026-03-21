package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"tgBookBot/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO CreateUser,CreateBook,GetUserBooks,DeleteBook,GetUserByTelegramID

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, dbCfg *config.DBConfig, log *slog.Logger) (pool *pgxpool.Pool, err error) {

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			log.Error("Attempt to connect failed, reconnecting", slog.Int("attempt", dbCfg.MaxAttempts), slog.Any("err", err))
			return err
		}
		return nil
	}, dbCfg.MaxAttempts, time.Second*dbCfg.DelayAttempts)
	if err != nil {
		log.Error("Failed to connect", slog.Any("err", err))
		return nil, err
	}
	return pool, nil

}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
