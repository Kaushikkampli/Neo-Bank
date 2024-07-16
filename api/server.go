package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/token"
	"github.com/kaushikkampli/neobank/utils"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	token  token.Maker
	config utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:  store,
		token:  token,
		config: config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupServer()

	return server, nil
}

func (server *Server) SetupServer() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRouter := router.Group("/", authMiddleware(server.token))

	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.ListAccounts)

	authRouter.POST("/transfer", server.CreateTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
