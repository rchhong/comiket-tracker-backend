package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Dbpool *pgxpool.Pool
}

func InitializeDB(host string, port int, databaseName string, username string, password string, poolMinConnections int, poolMaxConnections int) (*PostgresDB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_min_conns=%d&pool_max_conns=%d", username, password, host, port, databaseName, poolMinConnections, poolMaxConnections)
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	postgresDB := &PostgresDB{
		Dbpool: dbpool,
	}

	return postgresDB, nil
}
func (postgresdb *PostgresDB) Teardown() {
	postgresdb.Dbpool.Close()
}
