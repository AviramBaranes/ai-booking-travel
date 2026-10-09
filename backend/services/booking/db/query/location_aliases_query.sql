-- name: InsertManyLocationAliases :exec
INSERT INTO location_aliases (
    location_id,
    alias
)
SELECT
    location_ids.location_id,
    aliases.alias
FROM unnest(sqlc.arg(location_ids)::bigint[]) WITH ORDINALITY AS location_ids(location_id, idx)
JOIN unnest(sqlc.arg(aliases)::text[]) WITH ORDINALITY AS aliases(alias, idx)
    USING (idx)
ON CONFLICT (location_id, lower(alias))
DO NOTHING;

-- name: ListLocationAliasesByLocationIDs :many
SELECT location_id, alias
FROM location_aliases
WHERE location_id = ANY (sqlc.arg(location_ids)::bigint[])
ORDER BY location_id, lower(alias);

-- name: MoveLocationAliases :exec
INSERT INTO location_aliases (location_id, alias)
SELECT sqlc.arg(to_id)::bigint, alias
FROM location_aliases
WHERE location_id = sqlc.arg(from_id)::bigint
ON CONFLICT (location_id, lower(alias))
DO NOTHING;
