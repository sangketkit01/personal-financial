package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/util"
)

func (server *Server) GetFinancialById(ctx *gin.Context) {
	u, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized"))
		return
	}
	_, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("invalid user type"))
		return
	}

	financialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || financialId <= 0 {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid financial id."))
		return
	}

	financialData, err := server.store.GetFinancialById(ctx, int64(financialId))
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("no financial found."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial data."))
		return
	}

	ctx.JSON(http.StatusOK, financialData)
}

func (server *Server) MyFinancial(ctx *gin.Context) {
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

	myFinancial, err := server.store.MyFinancial(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial data."))
		return
	} else if len(myFinancial) == 0 {
		ctx.JSON(http.StatusNotFound, newErrorResponse("no financial found."))
		return
	}

	ctx.JSON(http.StatusOK, myFinancial)
}

type NewFinancialRequest struct {
	Amount int64  `json:"amount" binding:"required"`
	Type   string `json:"type" binding:"required,alpha"`
}

func (server *Server) AddNewFinancial(ctx *gin.Context) {
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

	var req NewFinancialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid request."))
		return
	}

	if req.Amount == 0 {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("amount cannot be zero"))
		return
	}

	financialTypeId, err := server.store.GetFinancialByName(ctx, util.CapitalizeWord(req.Type))
	if err != nil {
		if err == pgx.ErrNoRows {
			// other financial type
			financialTypeId.ID = 10
		} else {
			fmt.Printf("error: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial type."))
			return
		}
	}

	direction := "in"
	if req.Amount < 0 {
		direction = "out"
	}

	arg := db.InsertNewFinancialParams{
		UserID:    user.Username,
		Amount:    req.Amount,
		Direction: direction,
		TypeID:    financialTypeId.ID,
	}

	financial, err := server.store.InsertNewFinancial(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("failed to save your financial."))
		return
	}

	// after insert the new financial, check how much we've spent compared to our budget
	usageMessage := "you have no budget now. Please visit /insert-budget to add your budget"

	month := time.Now().Month()
	year := time.Now().Year()

	budget, err := server.store.GetBudget(ctx, db.GetBudgetParams{
		UserID: user.Username,
		Month:  int32(month),
		Year:   int32(year),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("user: %s have no budget\n", user.Username)
		} else {
			fmt.Printf("cannot get budget of the current user: %v", err)
		}
	} else {
		var usage float32 = 1.0

		currenMonthUsage, err := server.store.SummaryFinancialByMonth(ctx, db.SummaryFinancialByMonthParams{
			UserID: user.Username,
			Month:  int32(time.Now().Month()),
			Year:   int32(time.Now().Year()),
		})

		if err != nil {
			if err == pgx.ErrNoRows {
				usage = 0.0
			}
		}

		if budget.Amount.Int.Int64() == 0 {
			usageMessage = "Your budget is set to 0. Please update it to track usage."
		} else {
			budgetFloat, _ := budget.Amount.Float64Value()
			budgetAmount := float32(budgetFloat.Float64)
			usage = float32(math.Abs(float64(currenMonthUsage.TotalExpense))) / budgetAmount

			usagePercent := fmt.Sprintf("%.2f", usage*100.0)
			usageMessage = fmt.Sprintf("You've used %s%% of the budget you've set", usagePercent)
		}
	}
	// -----------------------------------------------------------------

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "saved financial successfully.",
		"financial": financial,
		"usage":     usageMessage,
	})
}

type UpdateFinancialRequest struct {
	Amount int64  `json:"amount" binding:"required"`
	Type   string `json:"type" binding:"required,alpha"`
}

func (server *Server) UpdateFinancial(ctx *gin.Context) {
	u, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized"))
		return
	}
	_, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("invalid user type"))
		return
	}

	financialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || financialId <= 0 {
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

	financialTypeId, err := server.store.GetFinancialByName(ctx, util.CapitalizeWord(req.Type))
	if err != nil {
		if err == pgx.ErrNoRows {
			// other financial type
			financialTypeId.ID = 10
		} else {
			ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get financial type."))
			return
		}
	}

	direction := "in"
	if req.Amount < 0 {
		direction = "out"
	}

	arg := db.UpdateFinancialParams{
		Amount:    req.Amount,
		Direction: direction,
		TypeID:    financialTypeId.ID,
		ID:        int64(financialId),
	}

	updatedFinancial, err := server.store.UpdateFinancial(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "cannot update financial.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":           "update financial successfully.",
		"updated_financial": updatedFinancial,
	})
}

func (server *Server) DeleteFinancial(ctx *gin.Context) {
	u, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized"))
		return
	}
	_, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("invalid user type"))
		return
	}

	financialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || financialId <= 0 {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid financial id."))
		return
	}

	deleteFinancial, err := server.store.DeleteFinancial(ctx, int64(financialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to delete financial")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":           "delete financial successfully.",
		"deleted_financial": deleteFinancial,
	})
}

func (server *Server) SummaryCurrentMonth(ctx *gin.Context) {
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

	arg := db.SummaryFinancialByMonthParams{
		UserID: user.Username,
		Month:  int32(time.Now().Month()),
		Year:   int32(time.Now().Year()),
	}

	summary, err := server.store.SummaryFinancialByMonth(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("you have no financial yet."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (server *Server) SummaryCurrentYear(ctx *gin.Context) {
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

	arg := db.SummaryFinancialByYearParams{
		UserID: user.Username,
		Year:   int32(time.Now().Year()),
	}

	summary, err := server.store.SummaryFinancialByYear(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("you have no financial yet."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

type YearMonthRequest struct {
	Year  int `json:"year" binding:"required,min=2020"`
	Month int `json:"month" binding:"required,min=1,max=12"`
}

type YearRequest struct {
	Year int `json:"year" binding:"required,min=2020"`
}

func (server *Server) SummaryByMonthYear(ctx *gin.Context) {
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
	var req YearMonthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid request body."))
		return
	}

	if req.Year > time.Now().Year() {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid year."))
		return
	}

	summary, err := server.store.SummaryFinancialByMonth(ctx, db.SummaryFinancialByMonthParams{
		UserID: user.Username,
		Month:  int32(req.Month),
		Year:   int32(req.Year),
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("you have no financial yet."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (server *Server) SummaryByYear(ctx *gin.Context) {
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

	var req YearRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid request body."))
		return
	}

	if req.Year > time.Now().Year() {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid year."))
		return
	}

	summary, err := server.store.SummaryFinancialByYear(ctx, db.SummaryFinancialByYearParams{
		UserID: user.Username,
		Year:   int32(req.Year),
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("you have no financial yet."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (server *Server) SummaryEachYear(ctx *gin.Context) {
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

	summary, err := server.store.SummaryFinancialEachYear(ctx, user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	} else if len(summary) == 0 {
		ctx.JSON(http.StatusNotFound, newErrorResponse("no financial found."))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (server *Server) SummaryTypeByMonthYear(ctx *gin.Context) {
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

	var req YearMonthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req.Month = int(time.Now().Month())
		req.Year = time.Now().Year()
	}

	if req.Year > time.Now().Year() {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid year."))
	}

	summary, err := server.store.SummaryByTypeMonth(ctx, db.SummaryByTypeMonthParams{
		UserID: user.Username,
		Month:  int32(req.Month),
		Year:   int32(req.Year),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	} else if len(summary) == 0 {
		ctx.JSON(http.StatusNotFound, newErrorResponse("you have no financial yet."))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (server *Server) SummaryTypeByYear(ctx *gin.Context) {
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

	var req YearRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req.Year = time.Now().Year()
	}

	if req.Year > time.Now().Year() {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid year."))
	}

	summary, err := server.store.SummaryByTypeYear(ctx, db.SummaryByTypeYearParams{
		UserID: user.Username,
		Year:   int32(req.Year),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
		return
	} else if len(summary) == 0 {
		ctx.JSON(http.StatusNotFound, newErrorResponse("you have no financial yet."))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}
