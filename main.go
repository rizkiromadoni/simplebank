package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiromadoni/simplebank/api"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	"github.com/rizkiromadoni/simplebank/gapi"
	pb "github.com/rizkiromadoni/simplebank/pb"
	"github.com/rizkiromadoni/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGRPCServer(config, store)
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal("cannot listen:", err)
	}

	log.Println("starting GRPC server on", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot serve:", err)
	}
}

func runServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddr)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
