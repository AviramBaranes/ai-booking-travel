-- name: InsertManyLocationAliases :exec
INSERT INTO location_aliases (
    location_id,
    name,
    type,
    language_code
)
SELECT
    location_ids.location_id,
    names.name,
    types.type_text::location_alias_type,
    language_codes.language_code
FROM unnest(sqlc.arg(location_ids)::bigint[]) WITH ORDINALITY AS location_ids(location_id, idx)
JOIN unnest(sqlc.arg(names)::text[]) WITH ORDINALITY AS names(name, idx)
    USING (idx)
JOIN unnest(sqlc.arg(types)::text[]) WITH ORDINALITY AS types(type_text, idx)
    USING (idx)
JOIN unnest(sqlc.arg(language_codes)::text[]) WITH ORDINALITY AS language_codes(language_code, idx)
    USING (idx)
ON CONFLICT (location_id, language_code, type, lower(name))
DO NOTHING;