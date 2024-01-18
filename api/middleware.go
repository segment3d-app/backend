package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/segment3d-app/segment3d-be/token"
)

const (
	authorizationHeaderKey    = "authorization"
	authorizationHeaderBearer = "Bearer"
	authorizationPayloadKey   = "autorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			error := errors.New("authorization header is empty")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(error))
			return
		}

		field := strings.Fields(authorizationHeader)
		if len(field) < 2 {
			error := errors.New("authorization header is not valid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(error))
			return
		}

		authorizationType := strings.ToLower(field[0])
		if authorizationType != strings.ToLower(authorizationHeaderBearer) {
			error := errors.New("authorization type is not valid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(error))
			return
		}

		accessToken := field[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
