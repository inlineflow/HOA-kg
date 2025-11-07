-- name: InsertFlatOwner :one
INSERT INTO flat_owner(flat_id, resident_id)
values($1, $2)
RETURNING *;
