package gapi

import (
	"context"

	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	proto "github.com/rizkiromadoni/simplebank/pb"
	"github.com/rizkiromadoni/simplebank/util"
	"github.com/rizkiromadoni/simplebank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(c context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

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

func validateCreateUserRequest(req *proto.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validator.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return
}
