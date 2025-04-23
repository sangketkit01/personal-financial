package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/token"
	"github.com/sangketkit01/personal-financial/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
	store db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store, tokenMaker token.Maker) (*Server, error){
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	server.setupRoute()
	server.router.Run(":5315")

	return server, nil
}

func (server *Server) setupRoute(){
	router := gin.Default()

	router.POST("/create-user",server.createUserAPI)
	router.POST("/login-user",server.LoginUserAPI)

	// authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	

	server.router = router
}
