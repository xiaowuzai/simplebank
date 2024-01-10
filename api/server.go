package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/token"
	"github.com/xiaowuzai/simplebank/util"
)

// Server 为服务提供 http 请求
type Server struct {
	store      db.Store //  使用接口，方便 mock
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config // 配置文件
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
	router := gin.Default()

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", valiadCurrency)
	}

	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

// Start 在指定的地址启动服务
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// errorResponse 返回错误
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
