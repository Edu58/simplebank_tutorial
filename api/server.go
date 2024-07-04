package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricalKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		router:     gin.Default(),
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
	server.router.POST("/user", server.createUser)
	server.router.POST("/user/login", server.loginUser)
	server.router.POST("/user/token/refresh", server.refreshUserToken)

	protectedRouted := server.router.Group("/").Use(authMiddleware(server.tokenMaker))

	server.router.POST("/user/:username", server.getUser)
	protectedRouted.POST("/accounts", server.createAccount)
	protectedRouted.GET("/accounts", server.listAccounts)
	protectedRouted.GET("/accounts/:id", server.getAccount)
	protectedRouted.POST("/transfer", server.createTransfer)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
