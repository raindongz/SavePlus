package apis

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/randongz/save_plus/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authentication middleware
// 1. get request header. 
// 2. check if informations in header are correct.
// 3. verify token, pass it to router.
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorizaton type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

const TraceID = "trace_id"

func setTraceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader(TraceID)
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx.Set("trace_id", traceID)
		ctx.Next()
	}
}
