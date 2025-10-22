-- name: CreateHouse :one
INSERT INTO house(address)
values($1)
RETURNING *;
