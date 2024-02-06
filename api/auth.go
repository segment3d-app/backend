package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
)

type registerUserRequest struct {
	Password string `json:"password" binding:"required,min=8,alphanum"`
	Email    string `json:"email" binding:"required,email"`
}

type registerUserResponse struct {
	AccessToken string        `json:"accessToken"`
	Message     string        `json:"message"`
	User        *UserResponse `json:"user"`
}

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerUserRequest true "User registration details"
// @Success 200 {object} registerUserResponse "User registration successful"
// @Router /auth/signup [post]
func (server *Server) signup(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(req.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := &registerUserResponse{
		User:        ReturnUserResponse(&user),
		AccessToken: accessToken,
		Message:     "registration successful",
	}
	ctx.JSON(http.StatusOK, response)
}


type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string       `json:"accessToken"`
	Message     string       `json:"message"`
	User        UserResponse `json:"user"`
}

// @Summary Login user
// @Description Login user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginUserRequest true "User login details"
// @Success 200 {object} loginUserResponse "User login successful"
// @Router /auth/signin [post]
func (server *Server) signin(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(req.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	loginUserData := &loginUserResponse{
		AccessToken: accessToken,
		User:        *ReturnUserResponse(&user),
		Message:     "login successful",
	}

	ctx.JSON(http.StatusOK, loginUserData)
}