package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/rizkiromadoni/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
)

func (s *Server) authorizeUser(c context.Context, accessibleRoles []string) (*token.Payload, error) {
	mtdt, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata")
	}

	values := mtdt.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("authorization header not found")
	}

	fields := strings.Fields(values[0])
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	if strings.ToLower(fields[0]) != authorizationTypeBearer {
		return nil, fmt.Errorf("invalid authorization type")
	}

	accessToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("user does not have permission")
	}

	return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}

	return false
}
