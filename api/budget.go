package api

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
)

type BudgetRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

// AddNewBudget adds budget for the current month - year
// Can't be add budget for the specified month - year
func (server *Server) AddNewBudget(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized."))
		return
	}

	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("invalid user type."))
		return
	}

	month := time.Now().Month()
	year := time.Now().Year()

	var req BudgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	amount := pgtype.Numeric{
		Int:   big.NewInt(req.Amount),
		Valid: true,
		Exp:   0,
	}

	budget, err := server.store.AddNewBudget(ctx, db.AddNewBudgetParams{
		UserID: user.Username,
		Month:  int32(month),
		Year:   int32(year),
		Amount: amount,
	})

	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			ctx.JSON(http.StatusBadRequest, newErrorResponse("you already have budget in the current month."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, budget)
}

// AddNewUpdateBudgetBudget updates budget for the current month - year
// Can't be update for the specified month - year
func (server *Server) UpdateBudget(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized."))
		return
	}

	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("invalid user type."))
		return
	}

	month := time.Now().Month()
	year := time.Now().Year()

	// check existence of the budget
	budget, err := server.store.GetBudget(ctx, db.GetBudgetParams{
		Month:  int32(month),
		Year:   int32(year),
		UserID: user.Username,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("budget not found."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req BudgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	amount := pgtype.Numeric{
		Int:   big.NewInt(req.Amount),
		Valid: true,
		Exp:   0,
	}

	updatedBudget, err := server.store.UpdateBudget(ctx, db.UpdateBudgetParams{
		Amount: amount,
		ID:     budget.ID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedBudget)
}

func (server *Server) GetCurrentBudget(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized."))
		return
	}

	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("invalid user type."))
		return
	}

	month := time.Now().Month()
	year := time.Now().Year()

	// check existence of the budget
	budget, err := server.store.GetBudget(ctx, db.GetBudgetParams{
		Month:  int32(month),
		Year:   int32(year),
		UserID: user.Username,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("budget not found."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, budget)
}

func (server *Server) GetHistoryBudget(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized."))
		return
	}

	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("invalid user type."))
		return
	}

	budgets, err := server.store.GetBudgetHistory(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("cannot get budget history."))
		return
	}

	if len(budgets) == 0 {
		ctx.JSON(http.StatusNotFound, newErrorResponse("budget not found."))
		return
	}

	ctx.JSON(http.StatusOK, budgets)
}

type BudgetHistoryRequest struct {
	Year int `json:"year" binding:"required,min=2000"`
}

func (server *Server) GetBudgetHistoryByYear(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized."))
		return
	}

	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("invalid user type."))
		return
	}

	var req BudgetHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Year > time.Now().Year() {
		ctx.JSON(http.StatusBadRequest, newErrorResponse("invalid year."))
		return
	}

	budget, err := server.store.GetBudgetHistoryByYear(ctx, db.GetBudgetHistoryByYearParams{
		UserID: user.Username,
		Year:   int32(req.Year),
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, newErrorResponse("budget not found."))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, budget)
}

type CheckBudgetUsageResponse struct {
	Budget       float64 `json:"budget"`
	Spent        int     `json:"spent"`
	UsagePercent string  `json:"usage"`
}

func (server *Server) CheckBudgetUsage(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("unauthorized"))
		return
	}

	user, ok := u.(db.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse("invalid user type"))
		return
	}

	response := CheckBudgetUsageResponse{
		Budget:       0,
		Spent:        0,
		UsagePercent: "0%",
	}

	budget, err := server.store.GetBudget(ctx, db.GetBudgetParams{
		Month:  int32(time.Now().Month()),
		Year:   int32(time.Now().Year()),
		UserID: user.Username,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusOK, response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("failed to get budget data"))
		return
	}

	summary, err := server.store.SummaryFinancialByMonth(ctx, db.SummaryFinancialByMonthParams{
		UserID: user.Username,
		Month:  int32(time.Now().Month()),
		Year:   int32(time.Now().Year()),
	})

	if err != nil && err != pgx.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse("failed to get financial summary"))
		return
	}

	budgetAmount, err := budget.Amount.Float64Value()
	if err == nil {
		response.Budget = budgetAmount.Float64
	} else {
		log.Printf("cannot parse budget amount: %v", err)
	}

	if err != pgx.ErrNoRows {
		response.Spent = int(math.Abs(float64(summary.TotalExpense)))
	}

	if response.Budget > 0 {
		usagePercent := math.Abs(float64(response.Spent) / float64(response.Budget) * 100)
		response.UsagePercent = fmt.Sprintf("%.2f%%", usagePercent)
	}

	ctx.JSON(http.StatusOK, response)
}
