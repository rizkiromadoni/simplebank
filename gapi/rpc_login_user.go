package gapi

import (
	"context"

	"github.com/jackc/pgx/v5"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	proto "github.com/rizkiromadoni/simplebank/pb"
	"github.com/rizkiromadoni/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(c context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	user, err := s.store.GetUser(c, req.GetUsername())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	err = util.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v", err)
	}

	session, err := s.store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpireAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	response := &proto.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpireAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpireAt),
	}
	return response, nil
}