package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rizkiromadoni/simplebank/db/sqlc"
	"github.com/rizkiromadoni/simplebank/token"
	"github.com/rizkiromadoni/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create paseto maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil

}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	})

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.POST("/tokens/renew_access", s.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccounts)
	authRoutes.POST("/accounts", s.createAccount)

	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
