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

type GenerateColmapEvent struct {
	AssetID string `json:"asset_id"`
	Url     string `json:"url"`
	Type    string `json:"type"`
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
// @Security BearerAuth
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

	asset, err = publishGenerateColmapEvent(server, ctx, &asset, &user)
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

func publishGenerateColmapEvent(server *Server, ginCtx *gin.Context, asset *db.Assets, user *db.Users) (db.Assets, error) {
	// generate message
	msg, err := json.Marshal(GenerateColmapEvent{
		AssetID: asset.ID.String(),
		Url:     asset.AssetUrl,
		Type:    asset.AssetType,
	})
	if err != nil {
		return *asset, err
	}

	err = server.rabbitmq.PublishEvent("generate_colmap", msg)
	if err != nil {
		return *asset, err
	}

	if asset.Status == "created" {
		arg := db.UpdateAssetStatusParams{
			Uid:    uuid.NullUUID{UUID: user.Uid, Valid: true},
			ID:     asset.ID,
			Status: "generating colmap",
		}

		newAsset, err := server.store.UpdateAssetStatus(ginCtx, arg)
		if err != nil {
			return *asset, err
		}

		return newAsset, nil
	}

	return *asset, nil
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

type getMyAssetsResponse struct {
	Message string          `json:"message"`
	Assets  []AssetResponse `json:"assets"`
}

// GetMyAssets
// @Summary Get my assets
// @Description Retrieves a list of my assets
// @Tags assets
// @Accept json
// @Produce json
// @Success 200 {object} getMyAssetsResponse "Success: Returns all assets."
// @Failure 500 {object} ErrorResponse "Error: Internal Server Error"
// @Security BearerAuth
// @Router /assets/me [get]
func (server *Server) getMyAssets(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	assets, err := server.store.GetMyAssets(ctx, uuid.NullUUID{UUID: user.Uid, Valid: true})
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
		formattedAsset = returnAssetResponse(&fAsset, &user)

		formattedAssets = append(formattedAssets, formattedAsset)
	}

	ctx.JSON(http.StatusOK, getMyAssetsResponse{Message: "all assets returned", Assets: formattedAssets})
}

type removeAssetRequest struct {
	ID string `uri:"id"`
}

type removeAssetResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

// RemoveAsset
// @Summary Remove my asset
// @Description Remove my asset
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Success 200 {object} removeAssetResponse "Asset removed successfully"
// @Security BearerAuth
// @Router /assets/{id} [delete]
func (server *Server) removeAsset(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req removeAssetRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.RemoveAssetParams{
		Uid: uuid.NullUUID{UUID: payload.Uid, Valid: true},
		ID:  uuid.MustParse(req.ID),
	}

	asset, err := server.store.RemoveAsset(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, removeAssetResponse{Message: "Asset removed successfully", Asset: returnAssetResponse(&asset, &user)})
}

type UpdatePointCloudUrlRequest struct {
	URL string `json:"url" binding:"required"`
}

type UpdatePointCloudUrlParam struct {
	ID string `uri:"id" binding:"required"`
}

type UpdatePointCloudUrlResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

type GenerateSplatEvent struct {
	AssetID string `json:"asset_id"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}

// UpdatePointCloudUrl updates the URL of a point cloud asset
// @Summary Update point cloud URL
// @Description Updates the URL for a specific point cloud asset based on the provided ID
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Param   request  body   UpdatePointCloudUrlRequest     true  "Update Point Cloud URL Request"
// @Success 200 {object} UpdatePointCloudUrlResponse "URL updated successfully"
// @Security BearerAuth
// @Router /assets/pointcloud/{id} [patch]
func (server *Server) updatePointCloudUrl(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdatePointCloudUrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var param UpdatePointCloudUrlParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdatePointCloudUrlParams{
		Uid:           uuid.NullUUID{UUID: payload.Uid, Valid: true},
		ID:            uuid.MustParse(param.ID),
		PointCloudUrl: sql.NullString{String: req.URL, Valid: true},
	}

	asset, err := server.store.UpdatePointCloudUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	asset, err = publishGenerateGaussianEvent(server, ctx, &asset, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := UpdatePointCloudUrlResponse{
		Message: "update success",
		Asset:   returnAssetResponse(&asset, &user),
	}

	ctx.JSON(http.StatusOK, res)
}

func publishGenerateGaussianEvent(server *Server, ginCtx *gin.Context, asset *db.Assets, user *db.Users) (db.Assets, error) {
	// generate message
	msg, err := json.Marshal(GenerateSplatEvent{
		AssetID: asset.ID.String(),
		Url:     asset.PointCloudUrl.String,
		Type:    asset.AssetType,
	})
	if err != nil {
		return *asset, err
	}

	err = server.rabbitmq.PublishEvent("generate_splat", msg)
	if err != nil {
		return *asset, err
	}

	if asset.Status == "generating colmap" {
		arg := db.UpdateAssetStatusParams{
			Uid:    uuid.NullUUID{UUID: user.Uid, Valid: true},
			ID:     asset.ID,
			Status: "generating splat",
		}

		newAsset, err := server.store.UpdateAssetStatus(ginCtx, arg)
		if err != nil {
			return *asset, nil
		}

		return newAsset, nil
	}

	return *asset, nil
}

type UpdateGaussianUrlRequest struct {
	URL string `json:"url" binding:"required"`
}

type UpdateGaussianUrlParam struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateGaussianUrlResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

// UpdateGaussianUrl updates the URL of a point cloud asset
// @Summary Update point cloud URL
// @Description Updates the URL for a specific gaussian asset based on the provided ID
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Param   request  body   UpdateGaussianUrlRequest     true  "Update Gaussian URL Request"
// @Success 200 {object} UpdateGaussianUrlResponse "URL updated successfully"
// @Security BearerAuth
// @Router /assets/gaussian/{id} [patch]
func (server *Server) updateGaussianUrl(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateGaussianUrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var param UpdateGaussianUrlParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	asset, err := server.store.GetAssetsById(ctx, uuid.MustParse(param.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdateGaussianUrlParams{
		Uid:         uuid.NullUUID{UUID: payload.Uid, Valid: true},
		ID:          uuid.MustParse(param.ID),
		GaussianUrl: sql.NullString{String: req.URL, Valid: true},
	}

	asset, err = server.store.UpdateGaussianUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	res := UpdateGaussianUrlResponse{
		Message: "update success",
		Asset:   returnAssetResponse(&asset, &user),
	}

	ctx.JSON(http.StatusOK, res)
}
