package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/docs"
	"github.com/segment3d-app/segment3d-be/token"
	"github.com/segment3d-app/segment3d-be/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{config: *config, store: store, tokenMaker: tokenMaker}
	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	authenticatedRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// configure swagger docs
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// auth api
	router.POST("/api/auth/signin", server.signin)
	router.POST("/api/auth/signup", server.signup)

	// user api
	authenticatedRouter.GET("/api/user", server.getUserData)
	authenticatedRouter.PATCH("/api/user", server.updateUser)
	authenticatedRouter.PATCH("/api/user/password", server.changeUserPassword)

	server.router = router
}

func getUserPayload(ctx *gin.Context) (*token.Payload, error) {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, fmt.Errorf("payload is missing")
	}
	userPayload, ok := payload.(*token.Payload)
	if !ok {
		return nil, fmt.Errorf("payload structure is not corrent")
	}

	return userPayload, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
