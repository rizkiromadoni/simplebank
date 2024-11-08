package gapi

import (
	"fmt"

	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	pb "github.com/rizkiromadoni/simplebank/pb"
	"github.com/rizkiromadoni/simplebank/token"
	"github.com/rizkiromadoni/simplebank/util"
	"github.com/rizkiromadoni/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create paseto maker: %w", err)
	}

	server := &Server{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
