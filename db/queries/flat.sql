-- name: InsertFlats :copyfrom
INSERT INTO flat(flat_number, house_id) VALUES($1, $2);

-- name: GetFlatsForHouse :many
select * from flat
where house_id = $1;

-- name: GetFlatWithNumber :one
select * from flat
where flat_number = $1
and house_id = $2;

-- name: FlatNumberExistsInHouse :one
select exists (
  select 1 from flat
  where flat_number = $1
  and house_id = $2
);
