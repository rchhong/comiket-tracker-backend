package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Dbpool *pgxpool.Pool
}

func InitializeDB() *PostgresDB {
	database_url := fmt.Sprintf("postgres://%s:%s@comiket-db:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	dbpool, err := pgxpool.New(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to setup database: %s", err)
	}

	return &PostgresDB{
		Dbpool: dbpool,
	}
}
func (postgresdb *PostgresDB) Teardown() {
	postgresdb.Dbpool.Close()
}
