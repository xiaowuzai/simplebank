package api

import (
	"github.com/gin-gonic/gin"
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

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

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
