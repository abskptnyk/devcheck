package check

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoCheck struct {
	URL    string
	pinger func(ctx context.Context, url string) error
}

func (c *MongoCheck) Name() string {
	return "MongoDB reachable"
}

func (c *MongoCheck) Run(ctx context.Context) Result {
	ping := c.pinger
	if ping == nil {
		ping = func(ctx context.Context, url string) error {
			client, err := mongo.Connect(options.Client().ApplyURI(url))
			if err != nil {
				return err
			}
			defer client.Disconnect(ctx)
			return client.Ping(ctx, nil)
		}
	}

	if err := ping(ctx, c.URL); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "cannot reach MongoDB",
			Fix:     "make sure MongoDB is running and MONGODB_URI is correct",
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "MongoDB is reachable",
	}
}
