package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
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

	return server, nil
}
