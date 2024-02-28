package gapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/util"
	"github.com/xiaowuzai/simplebank/worker"
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
