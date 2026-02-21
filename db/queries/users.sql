-- name: CreateUser :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByStripeCustomerID :one
SELECT * FROM users
WHERE stripe_customer_id = $1
LIMIT 1;

-- name: UpdateUserStripeCustomerID :one
UPDATE users
SET stripe_customer_id = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserSubscription :one
UPDATE users
SET
    stripe_subscription_id = $2,
    subscription_status    = $3
WHERE stripe_customer_id = $1
RETURNING *;

-- name: GetUserSubscriptionStatus :one
SELECT subscription_status FROM users
WHERE id = $1
LIMIT 1;
