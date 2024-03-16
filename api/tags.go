package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
)

type GetTagBySearchKeywordQuery struct {
	Keyword string `form:"keyword" binding:"required"`
}

type GetTagBySearchKeywordResponse struct {
	Message string `json:"message"`
	Tags []db.Tags `json:"tags"`
}

// @Summary Get tags by search keyword
// @Description Retrieves a list of tags that match the given search keyword
// @Tags tags
// @Accept  json
// @Produce  json
// @Param   keyword query string true "Search keyword"
// @Success 200 {object} GetTagBySearchKeywordResponse "Success"
// @Security BearerAuth
// @Router /tags/search [get]
func (server *Server) GetTagBySearchKeyword(ctx *gin.Context) {
	var req GetTagBySearchKeywordQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tags, err := server.store.GetTagsByKeyword(ctx, sql.NullString{String: req.Keyword, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := GetTagBySearchKeywordResponse {
		Message: "tags data retrived succesfully",
		Tags: tags,
	}
	ctx.JSON(http.StatusOK, res)
}