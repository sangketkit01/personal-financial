package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationHeaderType = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authMiddleware(tokenMaker token.Maker) gin.HandlerFunc{
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

		user, err := server.store.GetUser(ctx, payload.Username)
		if err != nil{
			if err == sql.ErrNoRows{
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error" : "user not found, then why are you here ?"})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error" : "cannot get user data."})
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func (server *Server) FinancialMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(db.User)

		financialId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil || financialId <= 0{
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error" : "invalid financial id."})
			return
		}

		owner, err := server.store.GetFinancialOwner(ctx, int64(financialId))
		if err != nil{
			if err == sql.ErrNoRows{
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error" : "no financial found."})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}

		if user.Username != owner{
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "you are not authorized to access this financial record",
			})

			return			
		}

		ctx.Next()
	}
}