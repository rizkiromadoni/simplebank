package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbSource = "postgres://postgres:postgres@localhost:5432/simplebank?sslmode=disable"

var testStore Store

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
