package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sangketkit01/personal-financial/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationHeaderType = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New("authorization header is not provided")})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) > 2{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error" : errors.New("invalid authorization header format")})
			return
		}

		authType := fields[0]
		if authType != authorizationHeaderType{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error" : errors.New(fmt.Sprintf("unsupported authorization type: %s",authType))})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil{
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error" : errors.New(err.Error())})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}