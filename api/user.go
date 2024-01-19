package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
)

type userResponse struct {
	ID                string    `json:"id"`
	Avatar            string    `json:"avatar"`
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PhoneNumber       string    `json:"phone_number"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

type registerUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password    string `json:"password" binding:"required,min=8,alphanum"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=12"`
}

type registerUserResponse struct {
	AccessToken string        `json:"access_token"`
	Message     string        `json:"message"`
	User        *userResponse `json:"user"`
}

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags user
// @Accept json
// @Produce json
// @Param request body registerUserRequest true "User registration details"
// @Success 200 {object} registerUserResponse "User registration successful"
// @Router /user/signup [post]
func (server *Server) registerUser(ctx *gin.Context) {
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
		Username:    req.Username,
		Password:    hashedPassword,
		Email:       req.Email,
		PhoneNumber: sql.NullString{String: req.PhoneNumber, Valid: true},
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

	accessToken, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := &registerUserResponse{
		User:        returnUserResponse(&user),
		AccessToken: accessToken,
		Message:     "registration successful",
	}
	ctx.JSON(http.StatusOK, response)
}

func returnUserResponse(user *db.User) *userResponse {
	return &userResponse{
		ID:                user.ID.String(),
		Avatar:            user.Avatar.String,
		Username:          user.Username,
		FullName:          user.FullName.String,
		Email:             user.Email,
		PhoneNumber:       user.PhoneNumber.String,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	Message     string       `json:"message"`
	User        userResponse `json:"user"`
}

// @Summary Login user
// @Description Login user with the provided credentials
// @Tags user
// @Accept json
// @Produce json
// @Param request body loginUserRequest true "User login details"
// @Success 200 {object} loginUserResponse "User login successful"
// @Router /user/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetAccountByUsername(ctx, req.Username)
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

	accessToken, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	loginUserData := &loginUserResponse{
		AccessToken: accessToken,
		User:        *returnUserResponse(&user),
		Message:     "login successful",
	}

	ctx.JSON(http.StatusOK, loginUserData)
}

type getUserByIdParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// @Summary Get user by ID
// @Description Retrieve user information based on the provided user ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} userResponse "User information retrieved successfully"
// @Security BearerAuth
// @Router /user/{id} [get]
func (server *Server) getUserById(ctx *gin.Context) {
	var req getUserByIdParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUserById(ctx, ID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, returnUserResponse(&user))
}

type updateUserParam struct {
	ID string `uri:"id" binding:"required"`
}

type updateUserRequest struct {
	Avatar      string `json:"avatar"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type updateUserResponse struct {
	Message string        `json:"message"`
	User    *userResponse `json:"user"`
}

// @Summary Update user information
// @Description Update user information based on the provided user ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Param request body updateUserRequest true "User update details"
// @Success 200 {object} updateUserResponse "User information updated successfully"
// @Security BearerAuth
// @Router /user/{id} [patch]
func (server *Server) updateUser(ctx *gin.Context) {
	var param updateUserParam
	var req updateUserRequest
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(param.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	avatar := req.Avatar
	fullName := req.FullName
	email := req.Email
	phoneNumber := req.PhoneNumber
	if len(avatar) == 0 {
		avatar = user.Avatar.String
	}
	if len(fullName) == 0 {
		fullName = user.FullName.String
	}
	if len(email) == 0 {
		email = user.Email
	}
	if len(phoneNumber) == 0 {
		phoneNumber = user.PhoneNumber.String
	}

	user, err = server.store.UpdateUser(ctx, db.UpdateUserParams{ID: id, Avatar: sql.NullString{Valid: true, String: avatar}, Email: email, PhoneNumber: sql.NullString{Valid: true, String: phoneNumber}, FullName: sql.NullString{Valid: true, String: fullName}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := &updateUserResponse{
		Message: "user information has been successfully updated",
		User:    returnUserResponse(&user),
	}

	ctx.JSON(http.StatusOK, response)
}

type changeUserPasswordParam struct {
	ID string `uri:"id" binding:"required"`
}

type changeUserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type changeUserPasswordResponse struct {
	Message string        `json:"message"`
	User    *userResponse `json:"user"`
}

// @Summary Change user password
// @Description Change user password based on the provided user ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Param request body changeUserPasswordRequest true "User password update details"
// @Success 200 {object} changeUserPasswordResponse "User password updated successfully"
// @Security BearerAuth
// @Router /user/{id}/password [patch]
func (server *Server) changeUserPassword(ctx *gin.Context) {
	var req changeUserPasswordRequest
	var param changeUserPasswordParam

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(param.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newHashedPassword, err := util.HashedPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err = server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{ID: id, Password: newHashedPassword})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := &changeUserPasswordResponse{
		Message: "user password has been successfully updated",
		User:    returnUserResponse(&user),
	}

	ctx.JSON(http.StatusOK, response)
}
