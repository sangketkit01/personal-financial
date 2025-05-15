package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
)

func (server *Server) GetFinancialById(ctx *gin.Context){
	_ = ctx.MustGet("user").(db.User)

	financialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || financialId <= 0{
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid financial id."))
		return
	}

	financialData, err := server.store.GetFinancialById(ctx, int64(financialId))
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, newErrorResponse("no financial found."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial data."))
		return
	}

	ctx.JSON(http.StatusOK, financialData)
}

func (server *Server) MyFinancial(ctx *gin.Context) {
	user := ctx.MustGet("user").(db.User)
	
	myFinancial, err := server.store.MyFinancial(ctx, user.Username)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, newErrorResponse("no financial found."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial data."))
		return
	}

	ctx.JSON(http.StatusOK, myFinancial)
}

type NewFinancialRequest struct {
	Amount int64  `json:"amount" binding:"required"`
	Type   string `json:"type" binding:"required,alpha"`
}

func (server *Server) AddNewFinancial(ctx *gin.Context) {
	user := ctx.MustGet("user").(db.User)

	var req NewFinancialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid request."))
		return
	}

	if req.Amount == 0 {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("amount cannot be zero"))
		return
	}

	financialTypeId, err := server.store.GetFinancialByName(ctx, req.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			// other financial type
			financialTypeId.ID = 10
		}else{
			ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial type."))
			return
		}
	}

	direction := "in"
	if req.Amount < 0 {
		direction = "out"
	}

	arg := db.InsertNewFinancialParams{
		UserID: user.Username,
		Amount: req.Amount,
		Direction: direction,
		TypeID: financialTypeId.ID,
	}

	financial ,err := server.store.InsertNewFinancial(ctx, arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("failed to save your financial."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message" : "saved financial successfully.",
		"financial" : financial,
	})
}

type UpdateFinancialRequest struct{
	Amount int64  `json:"amount" binding:"required"`
	Type   string `json:"type" binding:"required,alpha"`
}

func (server *Server) UpdateFinancial(ctx *gin.Context) {
	_ = ctx.MustGet("user").(db.User)
	
	financialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || financialId <= 0{
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid financial id."))
		return
	}

	var req UpdateFinancialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid request."))
		return
	}

	if req.Amount == 0 {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("amount cannot be zero"))
		return
	}

	financialTypeId, err := server.store.GetFinancialByName(ctx, req.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			// other financial type
			financialTypeId.ID = 10
		}else{
			ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial type."))
			return
		}
	}

	direction := "in"
	if req.Amount < 0 {
		direction = "out"
	}	

	arg := db.UpdateFinancialParams{
		Amount: req.Amount,
		Direction: direction,
		TypeID: financialTypeId.ID,
		ID: int64(financialId),
	}

	updatedFinancial, err := server.store.UpdateFinancial(ctx, arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, "cannot update financial.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message" : "update financial successfully.",
		"updated_financial" : updatedFinancial,
	})
}

func (server *Server) DeleteFinancial(ctx *gin.Context) {
	_ = ctx.MustGet("user").(db.User)

	financialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || financialId <= 0{
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid financial id."))
		return
	}

	deleteFinancial ,err := server.store.DeleteFinancial(ctx, int64(financialId))
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, "failed to delete financial")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message" : "delete financial successfully.",
		"deleted_financial" : deleteFinancial,
	})
}