package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/token"
	"github.com/sangketkit01/personal-financial/util"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct{
	Username string `json:"username" binding:"required"`
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required,min=10,max=10"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct{
	Username string `json:"username"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (server *Server) createUser(ctx *gin.Context){
	var req CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"hashed password error"})
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Name: req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx,arg)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "username or email already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
	}

	response := CreateUserResponse{
		Username: user.Username,
		Name: user.Name,
		Email: user.Email,
		Phone: user.Phone,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}

type LoginUserRequest struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserRespose struct{
	Username string `json:"username"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TokenID string `json:"token_id"`
	AccessToken string `json:"access_token"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (server *Server) LoginUser(ctx *gin.Context){
	var req LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil{
		if err == pgx.ErrNoRows{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "cannot get user"})
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil{
		if err == bcrypt.ErrMismatchedHashAndPassword{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"invalid password"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	payload, err := token.NewPayload(user.Username, 24 * time.Hour)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	accessToken , err := server.tokenMaker.CreateToken(user.Username, 24 * time.Hour)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	response := LoginUserRespose{
		Username: user.Username,
		Email: user.Email,
		Name: user.Name,
		Phone: user.Phone,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		TokenID: payload.ID.String(),
		AccessToken: accessToken,
		IssuedAt: payload.IssuedAt,
		ExpiredAt: payload.ExpiredAt,
	}
	
	ctx.JSON(http.StatusOK, response)
}