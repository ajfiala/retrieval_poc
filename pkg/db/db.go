package db 

import(
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"context"
	"sync"
)

var (
	pool *pgxpool.Pool 
	once sync.Once
)

func GetPool() *pgxpool.Pool {
	once.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("POSTGRES_CONN"))
		if err != nil {
			panic(err)
		}
	})
	return pool
}