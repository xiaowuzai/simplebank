package gapi

import (
	"fmt"

	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/pb"
	"github.com/xiaowuzai/simplebank/token"
	"github.com/xiaowuzai/simplebank/util"
)

// Server 为服务提供 http 请求
type Server struct {
	store      db.Store //  使用接口，方便 mock
	tokenMaker token.Maker
	config     util.Config // 配置文件
	pb.UnimplementedSimpleBankServer
}

// NewServer 创建一个 HTTP server 并提供所有 API 路由
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %s", err.Error())
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
