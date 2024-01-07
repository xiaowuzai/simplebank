package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
)

// Server 为服务提供 http 请求
type Server struct {
	store  db.Store //  使用接口，方便 mock
	router *gin.Engine
}

// NewServer 创建一个 HTTP server 并提供所有 API 路由
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", valiadCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start 在指定的地址启动服务
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// errorResponse 返回错误
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
