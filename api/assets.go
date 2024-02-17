package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
)

type AssetResponse struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Slug          string       `json:"slug"`
	AssetType     string       `json:"assetType"`
	Status        string       `json:"status"`
	ThumbnailUrl  string       `json:"thumbnailUrl"`
	AssetUrl      string       `json:"assetUrl"`
	PointCloudUrl string       `json:"pointCloudUrl"`
	GaussianUrl   string       `json:"gaussianUrl"`
	IsPrivate     bool         `json:"isPrivate"`
	Likes         int64        `json:"likes"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
	User          UserResponse `json:"user"`
}

func returnAssetResponse(asset *db.Assets, user *db.Users) AssetResponse {
	return AssetResponse{
		ID:            asset.ID.String(),
		Title:         asset.Title,
		Slug:          asset.Slug,
		AssetType:     asset.AssetType,
		Status:        asset.Status,
		ThumbnailUrl:  asset.ThumbnailUrl,
		AssetUrl:      asset.AssetUrl,
		PointCloudUrl: asset.PointCloudUrl.String,
		GaussianUrl:   asset.GaussianUrl.String,
		IsPrivate:     asset.IsPrivate,
		Likes:         int64(asset.Likes),
		CreatedAt:     asset.CreatedAt.String(),
		UpdatedAt:     asset.UpdatedAt.String(),
		User:          *ReturnUserResponse(user),
	}
}

type CreateAssetRequest struct {
	Title     string `json:"title" binding:"required"`
	IsPrivate *bool  `json:"isPrivate" binding:"required"`
	AssetUrl  string `json:"assetUrl" binding:"required"`
	AssetType string `json:"assetType" binding:"required,oneof=images video"`
}

type CreateAssetsResponse struct {
	Asset   AssetResponse `json:"asset"`
	Message string        `json:"message"`
}

type getThumbnailResponse struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}

// createAsset creates a new asset with provided details
// @Summary Create new asset
// @Description Creates a new asset based on the title, privacy setting, asset URL, and asset type provided in the request.
//
//	It also attempts to retrieve a thumbnail for the asset from the specified asset URL.
//
// @Tags assets
// @Accept json
// @Produce json
// @Param CreateAssetRequest body CreateAssetRequest true "Create Asset Request"
// @Success 202 {object} CreateAssetsResponse "Asset creation successful, returns created asset details along with a success message."
// @Router /assets [post]
func (server *Server) createAsset(ctx *gin.Context) {
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

	var req CreateAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	slug := util.GenerateBaseSlug(req.Title)
	pattern := slug + "%"
	existingSlugs, err := server.store.GetSlug(ctx, pattern)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(existingSlugs) > 0 {
		slug = slug + fmt.Sprintf("-%d", (len(existingSlugs)+1))
	}

	urlLink := strings.Replace(req.AssetUrl, "files", "thumbnail", -1)

	resp, err := http.Get(fmt.Sprintf("%s%s", server.config.StorageUrl, urlLink))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("expected JSON response, got: %s", contentType)))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response getThumbnailResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateAssetParams{
		Uid:          uuid.NullUUID{UUID: user.Uid, Valid: true},
		Title:        req.Title,
		Slug:         slug,
		Status:       "created",
		AssetUrl:     req.AssetUrl,
		AssetType:    req.AssetType,
		ThumbnailUrl: response.Url,
		IsPrivate:    false,
		Likes:        0,
	}

	if req.IsPrivate != nil {
		arg.IsPrivate = *req.IsPrivate
	}

	asset, err := server.store.CreateAsset(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := CreateAssetsResponse{
		Message: "generate splat from model",
		Asset:   returnAssetResponse(&asset, &user),
	}

	ctx.JSON(http.StatusAccepted, res)
}

type getAllAssetsResponse struct {
	Message string          `json:"message"`
	Assets  []AssetResponse `json:"assets"`
}

// GetAllAssets godoc
// @Summary Get all assets
// @Description Retrieves a list of all assets, including their associated user details.
// @Tags assets
// @Accept json
// @Produce json
// @Success 200 {object} getAllAssetsResponse "Success: Returns all assets."
// @Failure 500 {object} ErrorResponse "Error: Internal Server Error"
// @Router /assets [get]
func (server *Server) getAllAssets(ctx *gin.Context) {
	assets, err := server.store.GetAllAssets(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var formattedAssets []AssetResponse

	for _, asset := range assets {
		var formattedAsset AssetResponse
		fAsset := db.Assets{
			ID:            asset.ID,
			Title:         asset.Title,
			Slug:          asset.Slug,
			AssetType:     asset.AssetType,
			Status:        asset.Status,
			ThumbnailUrl:  asset.ThumbnailUrl,
			AssetUrl:      asset.AssetUrl,
			PointCloudUrl: asset.PointCloudUrl,
			GaussianUrl:   asset.GaussianUrl,
			IsPrivate:     asset.IsPrivate,
			Likes:         asset.Likes,
			CreatedAt:     asset.CreatedAt,
			UpdatedAt:     asset.UpdatedAt,
			Uid:           asset.Uid,
		}
		fUser := db.Users{
			Uid:    asset.Uid.UUID,
			Email:  asset.Email.String,
			Avatar: asset.Avatar,
			Name:   asset.Name,
		}
		formattedAsset = returnAssetResponse(&fAsset, &fUser)

		formattedAssets = append(formattedAssets, formattedAsset)
	}

	ctx.JSON(http.StatusOK, getAllAssetsResponse{Message: "all assets returned", Assets: formattedAssets})
}
