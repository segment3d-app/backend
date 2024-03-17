package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
)

type GetTagBySearchKeywordQuery struct {
	Keyword string `form:"keyword"`
	Limit   int    `form:"limit" binding:"required"`
}

type GetTagBySearchKeywordResponse struct {
	Message string    `json:"message"`
	Tags    []db.Tags `json:"tags"`
}

// GetTagBySearchKeyword
// @Summary Get tags by search keyword
// @Description Retrieves a list of tags that match the given search keyword
// @Tags tags
// @Accept  json
// @Produce  json
// @Param   keyword query string true "Search keyword"
// @Param   limit query int true "Limit for number of tags returned"
// @Success 200 {object} GetTagBySearchKeywordResponse "Success"
// @Security BearerAuth
// @Router /tags/search [get]
func (server *Server) GetTagBySearchKeyword(ctx *gin.Context) {
	var req GetTagBySearchKeywordQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetTagsByKeywordParams{
		Column1: sql.NullString{String: req.Keyword, Valid: true},
		Limit:   int64(req.Limit),
	}

	tags, err := server.store.GetTagsByKeyword(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := GetTagBySearchKeywordResponse{
		Message: "tags data retrived succesfully",
		Tags:    tags,
	}
	ctx.JSON(http.StatusOK, res)
}
