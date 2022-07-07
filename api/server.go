package api

import (
	db "github.com/Sanyam-Garg/simplebankgo/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/accounts/update", server.updateAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}
