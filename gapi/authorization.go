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

func (s *Server) authorizeUser(c context.Context) (*token.Payload, error) {
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

	return s.tokenMaker.VerifyToken(fields[1])
}
