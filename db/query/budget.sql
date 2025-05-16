-- name: AddNewBudget :one
INSERT INTO budgets
    (user_id, month, year, amount)
VALUES
    ($1, $2, $3, $4)
RETURNING *;

-- name: GetBudget :one
SELECT * FROM budgets
WHERE month = $1 AND year = $2 
AND user_id = $3 ;

-- name: GetBudgetHistory :many
SELECT * FROM budgets
WHERE user_id = $1
LIMIT 12;

-- name: GetBudgetHistoryByYear :one
SELECT * FROM budgets
WHERE user_id = $1 AND year = $2
LIMIT 12;

-- name: UpdateBudget :one
UPDATE budgets
SET amount = $1
WHERE id = $2
RETURNING *;