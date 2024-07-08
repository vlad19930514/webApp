package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/vlad19930514/webApp/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("user", server.createUser)
	router.GET("user/:id", server.getUser)
	router.PATCH("user", server.updateUser)
	server.router = router
	return server
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
