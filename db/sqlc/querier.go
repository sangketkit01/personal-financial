// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"
)

type Querier interface {
	AddNewBudget(ctx context.Context, arg AddNewBudgetParams) (Budget, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteFinancial(ctx context.Context, id int64) (Financial, error)
	GetBudget(ctx context.Context, arg GetBudgetParams) (Budget, error)
	GetBudgetHistory(ctx context.Context, userID string) ([]Budget, error)
	GetBudgetHistoryByYear(ctx context.Context, arg GetBudgetHistoryByYearParams) (Budget, error)
	GetFinancialById(ctx context.Context, id int64) (GetFinancialByIdRow, error)
	GetFinancialByName(ctx context.Context, type_ string) (FinancialType, error)
	GetFinancialOwner(ctx context.Context, id int64) (string, error)
	GetUser(ctx context.Context, username string) (User, error)
	InsertNewFinancial(ctx context.Context, arg InsertNewFinancialParams) (Financial, error)
	LoginUser(ctx context.Context, username string) (LoginUserRow, error)
	MyFinancial(ctx context.Context, userID string) ([]MyFinancialRow, error)
	SummaryByTypeMonth(ctx context.Context, arg SummaryByTypeMonthParams) ([]SummaryByTypeMonthRow, error)
	SummaryByTypeYear(ctx context.Context, arg SummaryByTypeYearParams) ([]SummaryByTypeYearRow, error)
	SummaryFinancialByMonth(ctx context.Context, arg SummaryFinancialByMonthParams) (SummaryFinancialByMonthRow, error)
	SummaryFinancialByYear(ctx context.Context, arg SummaryFinancialByYearParams) (SummaryFinancialByYearRow, error)
	SummaryFinancialEachYear(ctx context.Context, userID string) ([]SummaryFinancialEachYearRow, error)
	UpdateBudget(ctx context.Context, arg UpdateBudgetParams) (Budget, error)
	UpdateFinancial(ctx context.Context, arg UpdateFinancialParams) (Financial, error)
	UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error
}

var _ Querier = (*Queries)(nil)
