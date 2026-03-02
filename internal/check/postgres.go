package check

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresCheck struct {
	URL    string
	dialer func(ctx context.Context, url string) error
}

func (c *PostgresCheck) Name() string {
	return "PostgreSQL reachable"
}

func (c *PostgresCheck) Run(ctx context.Context) Result {
	dial := c.dialer
	if dial == nil {
		dial = func(ctx context.Context, url string) error {
			conn, err := pgx.Connect(ctx, url)
			if err != nil {
				return err
			}
			defer conn.Close(ctx)
			return conn.Ping(ctx)
		}
	}

	if err := dial(ctx, c.URL); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "cannot reach PostgreSQL",
			Fix:     "make sure PostgreSQL is running and DATABASE_URL is correct",
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "PostgreSQL is reachable",
	}
}
