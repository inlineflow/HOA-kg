-- name: InsertResident :one
INSERT INTO resident(flat_id, phone_number, fio)
VALUES($1, $2, $3)
RETURNING *;


