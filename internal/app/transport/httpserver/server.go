package httpserver

import "github.com/gin-gonic/gin"

type HttpServer struct {
	userService IUserService
	router      *gin.Engine
}

func NewHttpServer(userService IUserService) HttpServer {
	server := HttpServer{
		userService: userService,
	}

	router := gin.Default()
	router.POST("user", server.createUser)
	router.GET("user/:id", server.getUser)
	router.PATCH("user", server.updateUser)

	server.router = router
	return server
}

func (server *HttpServer) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
