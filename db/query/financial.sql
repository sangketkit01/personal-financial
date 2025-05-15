
-- name: InsertNewFinancial :one
INSERT INTO financials
    (user_id, amount, direction, type_id)
VALUES 
    ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateFinancial :one
UPDATE financials
SET amount = $1, direction = $2, type_id = $3
WHERE id = $4
RETURNING *;

-- name: DeleteFinancial :one
DELETE FROM financials 
WHERE id = $1
RETURNING *;

-- name: GetFinancialById :one
SELECT f.id, f.amount, f.direction, ft.type, f.created_at
FROM financials f
LEFT JOIN financial_types ft ON (f.type_id = ft.id) 
WHERE f.id = $1;

-- name: GetFinancialOwner :one
SELECT user_id FROM financials
WHERE id = $1;

-- name: MyFinancial :many
SELECT f.id, f.amount, f.direction, ft.type, f.created_at
FROM financials f
LEFT JOIN financial_types ft ON (f.type_id = ft.id) 
WHERE f.user_id = $1;