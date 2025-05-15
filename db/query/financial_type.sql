-- name: GetFinancialByName :one
SELECT * FROM financial_types
WHERE type = $1;