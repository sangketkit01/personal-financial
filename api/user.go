package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/token"
	"github.com/sangketkit01/personal-financial/util"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,min=10,max=10"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "hashed password error"})
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "username or email already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
	}

	response := CreateUserResponse{
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserRespose struct {
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TokenID     string    `json:"token_id"`
	AccessToken string    `json:"access_token"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get user"})
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payload, err := token.NewPayload(user.Username, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := LoginUserRespose{
		Username:    user.Username,
		Email:       user.Email,
		Name:        user.Name,
		Phone:       user.Phone,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		TokenID:     payload.ID.String(),
		AccessToken: accessToken,
		IssuedAt:    payload.IssuedAt,
		ExpiredAt:   payload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}

type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,alphanum"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqField=NewPassword"`
}

func (server *Server) UpdateUserPassword(ctx *gin.Context) {
	u, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized"))
		return
	}
	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("invalid user type"))
		return
	}

	var req UpdateUserPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)

	if req.CurrentPassword == req.NewPassword {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("new password must be different from current password"))
		return
	}

	if err := util.CheckPassword(req.CurrentPassword, user.Password); err != nil {
		ctx.JSON(http.StatusForbidden, newErrorResponse("invalid credentials"))
		return
	}

	newHashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.UpdatePassword(ctx, db.UpdatePasswordParams{
		Password: newHashedPassword,
		Username: user.Username,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("update password failed."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "update password successfully."})
}
