package db

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/pgvector/pgvector-go"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
)

var (
	pool *pgxpool.Pool 
	once sync.Once
)

func GetPool() *pgxpool.Pool {
	once.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("POSTGRES_CONN_STRING"))
		if err != nil {
			panic(err)
		}
	})
	return pool
}

func RegisterType() error {
	ctx := context.Background()
	// create *pgx.Conn
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		return err
	}

	err = pgxvector.RegisterTypes(ctx, conn)
	if err != nil {
		return err
	}
	return nil
}