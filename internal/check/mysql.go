package check

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLCheck struct {
	URL    string
	dialer func(ctx context.Context, url string) error
}

func (c *MySQLCheck) Name() string {
	return "MySQL reachable"
}

func (c *MySQLCheck) Run(ctx context.Context) Result {
	dial := c.dialer
	if dial == nil {
		dial = func(ctx context.Context, url string) error {
			db, err := sql.Open("mysql", url)
			if err != nil {
				return err
			}
			defer db.Close()
			return db.PingContext(ctx)
		}
	}

	if err := dial(ctx, c.URL); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "cannot reach MySQL",
			Fix:     "make sure MySQL is running and MYSQL_URL is correct",
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "MySQL is reachable",
	}
}
