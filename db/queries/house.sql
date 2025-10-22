-- name: CreateHouse :one
INSERT INTO house(address)
values($1)
RETURNING *;

-- name: GetAllHouses :many
SELECT * from house;
