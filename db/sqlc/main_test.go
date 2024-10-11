package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const dbSource = "postgres://postgres:postgres@localhost:5432/simplebank?sslmode=disable"

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
