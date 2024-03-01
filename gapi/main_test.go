package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/token"
	"github.com/xiaowuzai/simplebank/util"
	"github.com/xiaowuzai/simplebank/worker"
	"google.golang.org/grpc/metadata"
)

// newTestServer 创建一个测试服务器
func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey: util.NewPasetoSymmetricKey(),
		TokenDuration:     time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username, role string, timeDuration time.Duration) context.Context {
	token, _, err := tokenMaker.CreateToken(username, role, time.Minute)
	bearToken := fmt.Sprintf("%s %s", authorizationBearer, token)
	require.NoError(t, err)

	md := metadata.MD{
		authorizationHeader: []string{bearToken},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
