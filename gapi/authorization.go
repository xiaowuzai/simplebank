package gapi

import (
	"context"
	"errors"
	"strings"

	"github.com/xiaowuzai/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

var (
	ErrMissingMetadata           = errors.New("missing metadata")
	ErrAuthorizationHeader       = errors.New("missing authorization header")
	ErrAuthorizationHeaderFormat = errors.New("invaild authorization format")
	ErrAuthorizationType         = errors.New("invaild authorization type")
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingMetadata
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, ErrAuthorizationHeader
	}

	// 认证格式是：  bearer xxxxxx
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, ErrAuthorizationHeaderFormat
	}

	authType := strings.ToLower(fields[0])
	// switch authType {
	// case authorizationBearer:
	// default:
	// 	return nil, ErrAuthorizationHeaderFormat
	// }

	if authType != authorizationBearer {
		return nil, ErrAuthorizationType
	}

	accessToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}
	return payload, nil

}
