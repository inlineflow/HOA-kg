-- name: CreateFlats :copyfrom
INSERT INTO flat(flat_number, house_id) VALUES($1, $2);
