package gapi

import (
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	proto "github.com/rizkiromadoni/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *proto.User {
	return &proto.User{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
