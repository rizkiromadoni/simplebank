package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiromadoni/simplebank/api"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	"github.com/rizkiromadoni/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
