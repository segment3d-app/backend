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
	"github.com/lib/pq"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
)

type AssetResponse struct {
	ID                   string       `json:"id"`
	Title                string       `json:"title"`
	Slug                 string       `json:"slug"`
	Type                 string       `json:"type"`
	ThumbnailUrl         string       `json:"thumbnailUrl"`
	PhotoDirUrl          string       `json:"photoDirUrl"`
	SplatUrl             string       `json:"splatUrl"`
	PCLUrl               string       `json:"pclUrl"`
	PCLColmapUrl         string       `json:"pclColmapUrl"`
	SegmentedPclDirUrl   string       `json:"segmentedPclDirUrl"`
	SegmentedSplatDirUrl string       `json:"segmentedSplatDirUrl"`
	IsPrivate            bool         `json:"isPrivate"`
	Status               string       `json:"status"`
	Likes                int64        `json:"likes"`
	CreatedAt            string       `json:"createdAt"`
	UpdatedAt            string       `json:"updatedAt"`
	User                 UserResponse `json:"user"`
	IsLikedByMe          bool         `json:"isLikedByMe"`
}

type ReturnAssetResponseArg struct {
	Asset       *db.Assets
	User        *db.Users
	IsLikedByMe bool
}

func ReturnAssetResponse(arg ReturnAssetResponseArg) AssetResponse {
	return AssetResponse{
		ID:                   arg.Asset.ID.String(),
		Title:                arg.Asset.Title,
		Slug:                 arg.Asset.Slug,
		Type:                 arg.Asset.Type,
		ThumbnailUrl:         arg.Asset.ThumbnailUrl,
		PhotoDirUrl:          arg.Asset.PhotoDirUrl,
		SplatUrl:             arg.Asset.SplatUrl.String,
		PCLUrl:               arg.Asset.PclUrl.String,
		PCLColmapUrl:         arg.Asset.PclColmapUrl.String,
		SegmentedPclDirUrl:   arg.Asset.SegmentedPclDirUrl.String,
		SegmentedSplatDirUrl: arg.Asset.SegmentedSplatDirUrl.String,
		IsPrivate:            arg.Asset.IsPrivate,
		Likes:                int64(arg.Asset.Likes),
		Status:               arg.Asset.Status,
		CreatedAt:            arg.Asset.CreatedAt.String(),
		UpdatedAt:            arg.Asset.UpdatedAt.String(),
		User:                 *ReturnUserResponse(arg.User),
		IsLikedByMe:          arg.IsLikedByMe,
	}
}

type CreateAssetRequest struct {
	Title       string   `json:"title" binding:"required"`
	IsPrivate   *bool    `json:"isPrivate" binding:"required"`
	PhotoDirUrl string   `json:"photoDirUrl" binding:"required"`
	PCLUrl      string   `json:"pclUrl"`
	Type        string   `json:"type" binding:"required,oneof=lidar non_lidar"`
	Tags        []string `json:"tags"`
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
	AssetID       string `json:"asset_id"`
	PhotoDirUrl   string `json:"photo_dir_url"`
	Type          string `json:"type"`
	PointCloudUrl string `json:"point_cloud_url"`
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

	urlLink := strings.Replace(req.PhotoDirUrl, "files", "thumbnail", -1)

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
		Uid:          user.Uid,
		Title:        req.Title,
		Slug:         slug,
		Status:       "created",
		PhotoDirUrl:  req.PhotoDirUrl,
		Type:         req.Type,
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
	if len(req.PCLUrl) > 0 {
		argPCL := db.UpdatePointCloudUrlFromLidarParams{
			Uid:    user.Uid,
			ID:     asset.ID,
			PclUrl: sql.NullString{String: req.PCLUrl, Valid: true},
		}

		asset, err = server.store.UpdatePointCloudUrlFromLidar(ctx, argPCL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	tags, err := server.store.GetTagsByTagsName(ctx, req.Tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var allTags []db.Tags
	allTags = append(allTags, tags...)

	for _, curTag := range req.Tags {
		createTag := true
		for _, tag := range tags {
			if curTag == tag.Name {
				createTag = false
				break
			}
		}

		if createTag {
			tag, err := server.store.CreateTag(ctx, db.CreateTagParams{
				Name: curTag,
				Slug: util.GenerateBaseSlug(curTag),
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			allTags = append(allTags, tag)
		}
	}

	for _, tag := range allTags {
		_, err := server.store.CreateAssetsToTags(ctx, db.CreateAssetsToTagsParams{
			AssetsId: asset.ID,
			TagsId:   tag.ID,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	asset, err = publishGenerateColmapEvent(server, ctx, &asset, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := CreateAssetsResponse{
		Message: "generate splat from model",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user}),
	}

	ctx.JSON(http.StatusAccepted, res)
}

func publishGenerateColmapEvent(server *Server, ginCtx *gin.Context, asset *db.Assets, user *db.Users) (db.Assets, error) {
	// generate message
	msg, err := json.Marshal(GenerateColmapEvent{
		AssetID:       asset.ID.String(),
		PhotoDirUrl:   asset.PhotoDirUrl,
		Type:          asset.Type,
		PointCloudUrl: asset.PclUrl.String,
	})
	if err != nil {
		return *asset, err
	}

	err = server.rabbitmq.PublishEvent("process", msg)
	if err != nil {
		return *asset, err
	}

	if asset.Status == "created" {
		arg := db.UpdateAssetStatusParams{
			Uid:    user.Uid,
			ID:     asset.ID,
			Status: "generating sparse point cloud",
		}

		newAsset, err := server.store.UpdateAssetStatus(ginCtx, arg)
		if err != nil {
			return *asset, err
		}

		return newAsset, nil
	}

	return *asset, nil
}

type getAllAssetsQuery struct {
	Keyword string `form:"keyword"`
	Filter  string `form:"filter"`
}

type getAllAssetsResponse struct {
	Message string          `json:"message"`
	Assets  []AssetResponse `json:"assets"`
}

// GetAllAssets godoc
// @Summary Get all assets
// @Description Retrieves a list of all assets, optionally filtered by keyword and tags, including their associated user details.
// @Tags assets
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword for searching assets by title"
// @Param filter query string false "Comma-separated list of tags to filter the assets"
// @Success 200 {object} getAllAssetsResponse "Success: Returns all assets."
// @Failure 400 {object} ErrorResponse "Error: Bad Request"
// @Failure 500 {object} ErrorResponse "Error: Internal Server Error"
// @Router /assets [get]
func (server *Server) getAllAssets(ctx *gin.Context) {
	payload, errPayload := getUserPayload(ctx)

	var query getAllAssetsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	filterValues := strings.Split(query.Filter, ",")

	var formattedAssets []AssetResponse
	if errPayload != nil {
		assets, err := server.store.GetAllAssetsByKeyword(ctx, sql.NullString{String: query.Keyword, Valid: true})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var filteredAssets []db.GetAllAssetsByKeywordRow
		if len(query.Filter) > 0 {
			uniqueAssets := make(map[uuid.UUID]db.GetAllAssetsByKeywordRow)

			for _, asset := range assets {
				for _, curTag := range asset.TagNames {
					for _, tag := range filterValues {
						if curTag == tag {
							uniqueAssets[asset.ID] = asset
							break
						}
					}
				}
			}

			for _, asset := range uniqueAssets {
				filteredAssets = append(filteredAssets, asset)
			}
		} else {
			filteredAssets = assets
		}

		for _, asset := range filteredAssets {
			var formattedAsset AssetResponse
			fAsset := db.Assets{
				ID:                   asset.ID,
				Uid:                  asset.Uid,
				Title:                asset.Title,
				Slug:                 asset.Slug,
				Type:                 asset.Type,
				ThumbnailUrl:         asset.ThumbnailUrl,
				PhotoDirUrl:          asset.PhotoDirUrl,
				SplatUrl:             asset.SplatUrl,
				PclUrl:               asset.PclUrl,
				PclColmapUrl:         asset.PclColmapUrl,
				SegmentedPclDirUrl:   asset.SegmentedPclDirUrl,
				SegmentedSplatDirUrl: asset.SegmentedSplatDirUrl,
				IsPrivate:            asset.IsPrivate,
				Likes:                asset.Likes,
				Status:               asset.Status,
				CreatedAt:            asset.CreatedAt,
				UpdatedAt:            asset.UpdatedAt,
			}
			fUser := db.Users{
				Uid:    asset.Uid,
				Email:  asset.Email.String,
				Avatar: asset.Avatar,
				Name:   asset.Name,
			}
			formattedAsset = ReturnAssetResponse(ReturnAssetResponseArg{Asset: &fAsset, User: &fUser})
			formattedAssets = append(formattedAssets, formattedAsset)
		}
	} else {
		arg := db.GetAllAssetsWithLikesInformationParams{
			Uid:     payload.Uid,
			Column2: sql.NullString{String: query.Keyword, Valid: true},
		}
		assets, err := server.store.GetAllAssetsWithLikesInformation(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var filteredAssets []db.GetAllAssetsWithLikesInformationRow
		if len(query.Filter) > 0 {
			uniqueAssets := make(map[uuid.UUID]db.GetAllAssetsWithLikesInformationRow)

			for _, asset := range assets {
				for _, curTag := range asset.TagNames {
					for _, tag := range filterValues {
						if curTag == tag {
							uniqueAssets[asset.ID] = asset
							break
						}
					}
				}
			}

			for _, asset := range uniqueAssets {
				filteredAssets = append(filteredAssets, asset)
			}
		} else {
			filteredAssets = assets
		}

		for _, asset := range filteredAssets {
			var formattedAsset AssetResponse
			fAsset := db.Assets{
				ID:                   asset.ID,
				Uid:                  asset.Uid,
				Title:                asset.Title,
				Slug:                 asset.Slug,
				Type:                 asset.Type,
				ThumbnailUrl:         asset.ThumbnailUrl,
				PhotoDirUrl:          asset.PhotoDirUrl,
				SplatUrl:             asset.SplatUrl,
				PclUrl:               asset.PclUrl,
				PclColmapUrl:         asset.PclColmapUrl,
				SegmentedPclDirUrl:   asset.SegmentedPclDirUrl,
				SegmentedSplatDirUrl: asset.SegmentedSplatDirUrl,
				IsPrivate:            asset.IsPrivate,
				Likes:                asset.Likes,
				Status:               asset.Status,
				CreatedAt:            asset.CreatedAt,
				UpdatedAt:            asset.UpdatedAt,
			}
			fUser := db.Users{
				Uid:    asset.Uid,
				Email:  asset.Email.String,
				Avatar: asset.Avatar,
				Name:   asset.Name,
			}
			formattedAsset = ReturnAssetResponse(ReturnAssetResponseArg{Asset: &fAsset, User: &fUser, IsLikedByMe: asset.IsLikedByMe})
			formattedAssets = append(formattedAssets, formattedAsset)
		}
	}

	ctx.JSON(http.StatusOK, getAllAssetsResponse{Message: "all assets returned", Assets: formattedAssets})
}

type getMyAssetsQuery struct {
	Keyword string `form:"keyword"`
	Filter  string `form:"filter"`
}

type getMyAssetsResponse struct {
	Message string          `json:"message"`
	Assets  []AssetResponse `json:"assets"`
}

// GetMyAssets
// @Summary Get my assets
// @Description Retrieves a list of my assets, optionally filtered by keyword and tags.
// @Tags assets
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword for searching assets by title"
// @Param filter query string false "Comma-separated list of tags to filter the assets"
// @Success 200 {object} getMyAssetsResponse "Success: Returns all my assets."
// @Failure 400 {object} ErrorResponse "Error: Bad Request"
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

	var query getMyAssetsQuery
	if err = ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetMyAssetsParams{
		Uid: user.Uid,
		Column2: sql.NullString{
			String: query.Keyword,
			Valid:  true,
		},
	}

	assets, err := server.store.GetMyAssets(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	filterValues := strings.Split(query.Filter, ",")
	var filteredAssets []db.GetMyAssetsRow
	if len(query.Filter) > 0 {
		uniqueAssets := make(map[uuid.UUID]db.GetMyAssetsRow)

		for _, asset := range assets {
			for _, curTag := range asset.TagNames {
				for _, tag := range filterValues {
					if curTag == tag {
						uniqueAssets[asset.ID] = asset
						break
					}
				}
			}
		}

		for _, asset := range uniqueAssets {
			filteredAssets = append(filteredAssets, asset)
		}
	} else {
		filteredAssets = assets
	}

	var formattedAssets []AssetResponse

	for _, asset := range filteredAssets {
		var formattedAsset AssetResponse
		fAsset := db.Assets{
			ID:                   asset.ID,
			Uid:                  asset.Uid,
			Title:                asset.Title,
			Slug:                 asset.Slug,
			Type:                 asset.Type,
			ThumbnailUrl:         asset.ThumbnailUrl,
			PhotoDirUrl:          asset.PhotoDirUrl,
			SplatUrl:             asset.SplatUrl,
			PclUrl:               asset.PclUrl,
			PclColmapUrl:         asset.PclColmapUrl,
			SegmentedPclDirUrl:   asset.SegmentedPclDirUrl,
			SegmentedSplatDirUrl: asset.SegmentedSplatDirUrl,
			IsPrivate:            asset.IsPrivate,
			Likes:                asset.Likes,
			Status:               asset.Status,
			CreatedAt:            asset.CreatedAt,
			UpdatedAt:            asset.UpdatedAt,
		}
		formattedAsset = ReturnAssetResponse(ReturnAssetResponseArg{Asset: &fAsset, User: &user, IsLikedByMe: asset.IsLikedByMe.Bool})

		formattedAssets = append(formattedAssets, formattedAsset)
	}

	ctx.JSON(http.StatusOK, getMyAssetsResponse{Message: "all assets returned", Assets: formattedAssets})
}

// GetAssetDetails
// @Summary Get asset details
// @Description Get asset details by slug
// @Tags assets
// @Accept json
// @Produce json
// @Param slug path string true "Asset Slug"
// @Success 200 {object} getAssetDetailsResponse "Success response"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BearerAuth
// @Router /assets/{slug} [get]
type getAssetDetailsParams struct {
	Slug string `uri:"slug"`
}

type getAssetDetailsResponse struct {
	Asset   AssetResponse `json:"asset"`
	Message string        `json:"message"`
}

func (server *Server) getAssetDetails(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req getAssetDetailsParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	asset, err := server.store.GetAssetsBySlug(ctx, req.Slug)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	creator, err := server.store.GetUserById(ctx, asset.Uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	checkIsLikeArg := db.CheckIsLikedParams{
		Uid:      payload.Uid,
		AssetsId: asset.ID,
	}
	isLikedByMe, err := server.store.CheckIsLiked(ctx, checkIsLikeArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	assetResponseArg := ReturnAssetResponseArg{
		Asset:       &asset,
		User:        &creator,
		IsLikedByMe: isLikedByMe,
	}

	res := getAssetDetailsResponse{
		Asset:   ReturnAssetResponse(assetResponseArg),
		Message: "success to get assetDetails",
	}

	ctx.JSON(http.StatusOK, res)
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
		Uid: payload.Uid,
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

	ctx.JSON(http.StatusAccepted, removeAssetResponse{Message: "Asset removed successfully", Asset: ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user})})
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
	AssetID      string `json:"asset_id"`
	PhotoDirUrl  string `json:"photo_dir_url"`
	PCLColmapUrl string `json:"pcl_colmap_url"`
	Type         string `json:"type"`
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
// @Router /assets/pointcloud/{id} [patch]
func (server *Server) updatePointCloudUrl(ctx *gin.Context) {
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

	arg := db.UpdatePointCloudUrlFromColmapParams{
		ID:           uuid.MustParse(param.ID),
		PclColmapUrl: sql.NullString{String: req.URL, Valid: true},
	}

	asset, err := server.store.UpdatePointCloudUrlFromColmap(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, asset.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if asset.Status == "generating sparse point cloud" {
		arg := db.UpdateAssetStatusParams{
			Uid:    user.Uid,
			ID:     asset.ID,
			Status: "generating 3d splat",
		}

		asset, err = server.store.UpdateAssetStatus(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res := UpdatePointCloudUrlResponse{
		Message: "update success",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user}),
	}

	ctx.JSON(http.StatusOK, res)
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
// @Router /assets/gaussian/{id} [patch]
func (server *Server) updateGaussianUrl(ctx *gin.Context) {
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

	asset, err := server.store.GetAssetsById(ctx, uuid.MustParse(param.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdateSplatUrlParams{
		ID:       uuid.MustParse(param.ID),
		SplatUrl: sql.NullString{String: req.URL, Valid: true},
	}

	asset, err = server.store.UpdateSplatUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, asset.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if asset.Status == "generating 3d splat" {
		arg := db.UpdateAssetStatusParams{
			Uid:    user.Uid,
			ID:     asset.ID,
			Status: "processing ptv3",
		}

		asset, err = server.store.UpdateAssetStatus(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res := UpdateGaussianUrlResponse{
		Message: "update success",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user}),
	}

	ctx.JSON(http.StatusOK, res)
}

type UpdatePTV3UrlRequest struct {
	URL string `json:"url" binding:"required"`
}

type UpdatePTV3UrlParam struct {
	ID string `uri:"id" binding:"required"`
}

type UpdatePTV3UrlResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

// UpdatePTV3Url updates the URL of a PTv3 asset
// @Summary Update PTv3 URL
// @Description Updates the URL for a specific PTv3 asset based on the provided ID
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Param   request  body   UpdatePTV3UrlRequest     true  "Update PTv3 URL Request"
// @Success 200 {object} UpdatePTV3UrlResponse "URL updated successfully"
// @Router /assets/ptv3/{id} [patch]
func (server *Server) updatePTv3Url(ctx *gin.Context) {
	var req UpdatePTV3UrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var param UpdatePTV3UrlParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	asset, err := server.store.GetAssetsById(ctx, uuid.MustParse(param.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdatePTvUrlParams{
		ID:                 uuid.MustParse(param.ID),
		SegmentedPclDirUrl: sql.NullString{String: req.URL, Valid: true},
	}

	asset, err = server.store.UpdatePTvUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, asset.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if asset.Status == "processing ptv3" {
		arg := db.UpdateAssetStatusParams{
			Uid:    user.Uid,
			ID:     asset.ID,
			Status: "processing saga",
		}

		asset, err = server.store.UpdateAssetStatus(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res := UpdatePTV3UrlResponse{
		Message: "update success",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user}),
	}

	ctx.JSON(http.StatusOK, res)
}

type UpdateSagaUrlRequest struct {
	URL string `json:"url" binding:"required"`
}

type UpdateSagaUrlParam struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateSagaUrlResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

// UpdateSagaUrl updates the URL of a saga asset
// @Summary Update saga URL
// @Description Updates the URL for a specific saga asset based on the provided ID
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Param   request  body   UpdateSagaUrlRequest     true  "Update Saga URL Request"
// @Success 200 {object} UpdateSagaUrlResponse "URL updated successfully"
// @Router /assets/saga/{id} [patch]
func (server *Server) updateSagaUrl(ctx *gin.Context) {
	var req UpdateSagaUrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var param UpdateSagaUrlParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	asset, err := server.store.GetAssetsById(ctx, uuid.MustParse(param.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.UpdateSagaUrlParams{
		ID:                   uuid.MustParse(param.ID),
		SegmentedSplatDirUrl: sql.NullString{String: req.URL, Valid: true},
	}

	asset, err = server.store.UpdateSagaUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, asset.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if asset.Status == "processing saga" {
		arg := db.UpdateAssetStatusParams{
			Uid:    user.Uid,
			ID:     asset.ID,
			Status: "completed",
		}

		asset, err = server.store.UpdateAssetStatus(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res := UpdateSagaUrlResponse{
		Message: "update success",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user}),
	}

	ctx.JSON(http.StatusOK, res)
}

type LikeAssetParam struct {
	ID string `uri:"id" binding:"required"`
}

type LikeAssetResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

// LikeAsset handler to like an asset
// @Summary Like an asset
// @Description Marks an asset as liked by the current user.
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Security BearerAuth
// @Success 200 {object} LikeAssetResponse "Asset liked successfully"
// @Router /assets/like/{id} [post]
func (server *Server) likeAsset(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req LikeAssetParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	argCreateLike := db.CreateLikeParams{
		Uid:      payload.Uid,
		AssetsId: uuid.MustParse(req.ID),
	}
	err = server.store.CreateLike(ctx, argCreateLike)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			asset, err := server.store.GetAssetsById(ctx, uuid.MustParse(req.ID))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			res := LikeAssetResponse{
				Message: "asset already liked before",
				Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user, IsLikedByMe: true}),
			}
			ctx.JSON(http.StatusConflict, res)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	asset, err := server.store.IncreaseAssetLikes(ctx, uuid.MustParse(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := LikeAssetResponse{
		Message: "like asset success",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user, IsLikedByMe: true}),
	}

	ctx.JSON(http.StatusOK, res)
}

type UnlikeAssetParam struct {
	ID string `uri:"id" binding:"required"`
}

type UnlikeAssetResponse struct {
	Message string        `json:"message"`
	Asset   AssetResponse `json:"asset"`
}

// UnlikeAsset handler to unlike an asset
// @Summary Unlike an asset
// @Description Marks an asset as unliked by the current user, removing the like.
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Security BearerAuth
// @Success 200 {object} UnlikeAssetResponse "Asset unliked successfully"
// @Router /assets/unlike/{id} [post]
func (server *Server) unlikeAsset(ctx *gin.Context) {
	payload, err := getUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, payload.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var req UnlikeAssetParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.RemoveLikeParams{
		Uid:      payload.Uid,
		AssetsId: uuid.MustParse(req.ID),
	}
	_, err = server.store.RemoveLike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	asset, err := server.store.DecreaseAssetLikes(ctx, uuid.MustParse(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := UnlikeAssetResponse{
		Message: "unlike asset success",
		Asset:   ReturnAssetResponse(ReturnAssetResponseArg{Asset: &asset, User: &user}),
	}

	ctx.JSON(http.StatusOK, res)
}

type SegmentUsingSagaRequest struct {
	X                int    `json:"x" binding:"required"`
	Y                int    `json:"y" binding:"required"`
	URL              string `json:"url" binding:"required"`
	UniqueIdentifier string `json:"uniqueIdentifier" binding:"required"`
}

type SegmentUsingSagaParam struct {
	ID string `uri:"id" binding:"required"`
}

type SegmentUsingSagaResponse struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}

// SegmentUsingSaga Segment using SAGA
// @Summary Segment using SAGA
// @Description Segment using SAGA by sending message to RabbitMQ
// @Tags assets
// @Accept json
// @Produce json
// @Param   id   path   string     true  "Asset ID"
// @Param   request  body   SegmentUsingSagaRequest     true  "Segment using SAGA Request"
// @Success 200 {object} SegmentUsingSagaResponse "Segment using SAGA successfully"
// @Router /assets/saga/segment/{id} [post]
func (server *Server) segmentUsingSaga(ctx *gin.Context) {
	var req SegmentUsingSagaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var param SegmentUsingSagaParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	asset, err := server.store.GetAssetsById(ctx, uuid.MustParse(param.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = publishSegmentUsingSagaEvent(server, GenerateSegmentUsingSagaEvent{
		AssetID: asset.ID.String(),
		X: req.X,
		Y: req.Y,
		URL: req.URL,
		UniqueIdentifier: req.UniqueIdentifier,
	})

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	res := SegmentUsingSagaResponse {
		Message: "success",
		Url: fmt.Sprintf("/files/%s/%s.ply", asset.ID.String(), req.UniqueIdentifier),
	}

	ctx.JSON(http.StatusOK, res)
}

type GenerateSegmentUsingSagaEvent struct {
	AssetID          string `json:"asset_id"`
	X                int    `json:"x"`
	Y                int    `json:"y"`
	UniqueIdentifier string `json:"unique_identifier"`
	URL              string `json:"url"`
}

func publishSegmentUsingSagaEvent(server *Server, message GenerateSegmentUsingSagaEvent) error {
	// generate message
	msg, err := json.Marshal(message)

	if err != nil {
		return err
	}

	err = server.rabbitmq.PublishEvent("query", msg)
	if err != nil {
		return err
	}

	return nil
}
