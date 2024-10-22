package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiromadoni/simplebank/api"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
)

const (
	dbSource      = "postgres://postgres:postgres@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:3000"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
