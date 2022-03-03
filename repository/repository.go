package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
)

type DBLogger struct{}

func (logger DBLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (logger DBLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	fmt.Println(string(query))
	return nil
}

func New(opts *pg.Options) *pg.DB {
	return pg.Connect(opts)
}
