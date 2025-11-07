-- name: InsertFlat :one
insert into flat(flat_number, house_id, area)
values ($1, $2, $3)
RETURNING *;

-- name: GetFlatsForHouse :many
select sqlc.embed(f), sqlc.embed(r) from flat f
inner join flat_owner fo on f.flat_id = fo.flat_id
inner join resident r on fo.resident_id = r.resident_id
where house_id = $1
and fo.owner_up_to is null;

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
