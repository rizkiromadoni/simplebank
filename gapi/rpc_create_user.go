package gapi

import (
	"context"

	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	proto "github.com/rizkiromadoni/simplebank/pb"
	"github.com/rizkiromadoni/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(c context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := db.CreateUserParams{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		FullName: req.GetFullName(),
		Password: hashedPassword,
	}

	user, err := s.store.CreateUser(c, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	response := &proto.CreateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}
