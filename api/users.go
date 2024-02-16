package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
)

type UserResponse struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	Avatar            string    `json:"avatar"`
	Provider          string    `json:"provider"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
}

func ReturnUserResponse(user *db.Users) *UserResponse {
	return &UserResponse{
		ID:                user.Uid.String(),
		Avatar:            user.Avatar.String,
		Name:              user.Name.String,
		Email:             user.Email,
		Provider:          user.Provider,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

// @Summary Get user data
// @Description Retrieve user information
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse "User information retrieved successfully"
// @Security BearerAuth
// @Router /users [get]
func (server *Server) getUserData(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user with email %s is not found", user.Email)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ReturnUserResponse(&user))
}

type updateUserRequest struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

type updateUserResponse struct {
	Message string        `json:"message"`
	User    *UserResponse `json:"user"`
}

// @Summary Update user information
// @Description Update user information based on the provided user ID
// @Tags users
// @Accept json
// @Produce json
// @Param request body updateUserRequest true "User update details"
// @Success 200 {object} updateUserResponse "User information updated successfully"
// @Security BearerAuth
// @Router /users [patch]
func (server *Server) updateUser(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user with email %s is not found", user.Email)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	avatar := req.Avatar
	name := req.Name
	if len(avatar) == 0 {
		avatar = user.Avatar.String
	}
	if len(name) == 0 {
		name = user.Name.String
	}

	user, err = server.store.UpdateUser(ctx, db.UpdateUserParams{Uid: user.Uid, Avatar: sql.NullString{Valid: true, String: avatar}, Email: user.Email, Name: sql.NullString{Valid: true, String: name}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := &updateUserResponse{
		Message: "user information has been successfully updated",
		User:    ReturnUserResponse(&user),
	}

	ctx.JSON(http.StatusOK, response)
}

type changeUserPasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type changeUserPasswordResponse struct {
	Message string        `json:"message"`
	User    *UserResponse `json:"user"`
}

// @Summary Change user password
// @Description Change user password based on the provided user ID
// @Tags users
// @Accept json
// @Produce json
// @Param request body changeUserPasswordRequest true "User password update details"
// @Success 200 {object} changeUserPasswordResponse "User password updated successfully"
// @Security BearerAuth
// @Router /users/password [patch]
func (server *Server) changeUserPassword(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req changeUserPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user with email %s is not found", user.Email)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.OldPassword, user.Password.String)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("old password doesn't match")))
		return
	}

	newHashedPassword, err := util.HashedPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err = server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{Uid: user.Uid, Password: sql.NullString{String: newHashedPassword, Valid: true}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := &changeUserPasswordResponse{
		Message: "user password has been successfully updated",
		User:    ReturnUserResponse(&user),
	}

	ctx.JSON(http.StatusOK, response)
}
