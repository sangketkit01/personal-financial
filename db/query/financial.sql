
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

-- name: SummaryFinancialByMonth :one
SELECT 
  SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) AS total_income,
  SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) AS total_expense,
  CASE
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) > SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'in'
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) < SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'out'
    ELSE 'equal'
  END AS status
FROM financials f
WHERE f.user_id = @user_id::text
  AND EXTRACT(MONTH FROM f.created_at) = @month::int
  AND EXTRACT(YEAR FROM f.created_at) = @year::int;

-- name: SummaryFinancialByYear :one
SELECT 
  SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) AS total_income,
  SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) AS total_expense,
  CASE
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) > SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'in'
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) < SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'out'
    ELSE 'equal'
  END AS status
FROM financials f
WHERE f.user_id = @user_id::text
  AND EXTRACT(YEAR FROM f.created_at) = @year::int;


-- name: SummaryByTypeMonth :many
SELECT 
  ft.type,
  SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) AS total_income,
  SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) AS total_expense,
  CASE
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) > SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'in'
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) < SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'out'
    ELSE 'equal'
  END AS status
FROM financials f
JOIN financial_types ft ON ft.id = f.type_id
WHERE f.user_id = @user_id::text
  AND EXTRACT(MONTH FROM f.created_at) = @month::int
  AND EXTRACT(YEAR FROM f.created_at) = @year::int
GROUP BY ft.type
ORDER BY ft.type;

-- name: SummaryByTypeYear :many
SELECT 
  ft.type,
  SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) AS total_income,
  SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) AS total_expense,
  CASE
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) > SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'in'
    WHEN SUM(CASE WHEN f.direction = 'in' THEN f.amount ELSE 0 END) < SUM(CASE WHEN f.direction = 'out' THEN f.amount ELSE 0 END) THEN 'out'
    ELSE 'equal'
  END AS status
FROM financials f
JOIN financial_types ft ON ft.id = f.type_id
WHERE f.user_id = @user_id::text
  AND EXTRACT(YEAR FROM f.created_at) = @year::int
GROUP BY ft.type
ORDER BY ft.type;

-- name: SummaryFinancialEachYear :many
SELECT 
    EXTRACT(YEAR FROM f.created_at)::INT AS year,
    COALESCE(SUM(CASE WHEN f.direction = 'in' THEN f.amount END), 0) AS in_amount,
    COALESCE(SUM(CASE WHEN f.direction = 'out' THEN f.amount END), 0) AS out_amount,
    CASE
        WHEN COALESCE(SUM(CASE WHEN f.direction = 'in' THEN f.amount END), 0) > COALESCE(SUM(CASE WHEN f.direction = 'out' THEN f.amount END), 0) THEN 'in'
        WHEN COALESCE(SUM(CASE WHEN f.direction = 'in' THEN f.amount END), 0) < COALESCE(SUM(CASE WHEN f.direction = 'out' THEN f.amount END), 0) THEN 'out'
        ELSE 'equal'
    END AS status
FROM financials f
WHERE f.user_id = @user_id::text
GROUP BY year
ORDER BY year;
