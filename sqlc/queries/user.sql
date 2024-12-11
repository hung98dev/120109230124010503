-- sqlc/queries/user.sql
-- name: CreateUser :one
INSERT INTO public.user (
    site,
    employee_id,
    username,
    password_hash,
    department_id,
    role,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING employee_id;

-- name: GetUserByUserName :one
SELECT 
    id,
    site,
    employee_id,
    username,
    password_hash,
    email,
    phone,
    department_id,
    role,
    status,
    status_reason
FROM public.user
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmployeeID :one
SELECT 
    id,
    site,
    employee_id,
    username,
    password_hash,
    email,
    phone,
    department_id,
    role,
    status,
    status_reason FROM public.user
WHERE employee_id = $1 LIMIT 1;