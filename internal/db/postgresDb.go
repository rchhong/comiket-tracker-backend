package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Dbpool *pgxpool.Pool
}

func InitializeDB(host string, port int, databaseName string, username string, password string) (*PostgresDB, error) {
	database_url := fmt.Sprintf("postgres://%s:%s@comiket-db:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	dbpool, err := pgxpool.New(context.Background(), database_url)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		Dbpool: dbpool,
	}, nil
}
func (postgresdb *PostgresDB) Teardown() {
	postgresdb.Dbpool.Close()
}
