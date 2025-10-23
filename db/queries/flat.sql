-- name: CreateFlats :copyfrom
INSERT INTO flat(flat_number, house_id) VALUES($1, $2);

-- name: GetFlatsForHouse :many
select * from flat
where house_id = $1;
