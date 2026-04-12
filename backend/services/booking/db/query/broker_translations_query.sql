-- name: GetAllVerifiedTranslations :many
SELECT
    id,
    source_text,
    target_text
FROM
    broker_translations
WHERE
    status = 'verified';

-- name: GetAllTranslationSourceTexts :many
SELECT
    source_text
FROM
    broker_translations;

-- name: UpdateBrokerTranslation :exec
UPDATE broker_translations
SET
    target_text = sqlc.arg (target_text),
    status = 'verified',
    confidence_score = 10,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg (id);

-- name: VerifyBrokerTranslation :exec
UPDATE broker_translations
SET
    status = 'verified',
    confidence_score = 10,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg (id);

-- name: CheckBrokerTranslationExists :one
SELECT
    id
FROM
    broker_translations
WHERE
    source_text = sqlc.arg (source_text)
LIMIT
    1;

-- name: InsertBrokerTranslation :one
INSERT INTO
    broker_translations (source_text)
VALUES
    (sqlc.arg (source_text)) RETURNING id;

-- name: InsertBrokerTranslationFull :one
INSERT INTO
    broker_translations (source_text, target_text, status, confidence_score)
VALUES
    (sqlc.arg (source_text), sqlc.narg (target_text), sqlc.arg (status)::broker_translation_status, sqlc.arg (confidence_score)) RETURNING id;

-- name: DeleteBrokerTranslation :exec
DELETE FROM broker_translations
WHERE
    id = sqlc.arg (id);

-- name: ListAllTranslations :many
SELECT
    *
FROM
    broker_translations
WHERE
    (
        sqlc.narg (search)::text IS NULL
        OR source_text ILIKE '%' || sqlc.narg (search) || '%'
        OR target_text ILIKE '%' || sqlc.narg (search) || '%'
    )
    AND (
        sqlc.narg (status)::broker_translation_status IS NULL
        OR status = sqlc.narg (status)::broker_translation_status
    )
ORDER BY
    CASE
        WHEN sqlc.arg (sort_dir)::text = 'asc' THEN confidence_score
    END ASC,
    CASE
        WHEN sqlc.arg (sort_dir)::text = 'desc' THEN confidence_score
    END DESC,
    id ASC
LIMIT
    sqlc.arg (query_limit)::int
OFFSET
    sqlc.arg (query_offset)::int;

-- name: CountAllTranslations :one
SELECT
    COUNT(*)
FROM
    broker_translations
WHERE
    (
        sqlc.narg (search)::text IS NULL
        OR source_text ILIKE '%' || sqlc.narg (search) || '%'
        OR target_text ILIKE '%' || sqlc.narg (search) || '%'
    )
    AND (
        sqlc.narg (status)::broker_translation_status IS NULL
        OR status = sqlc.narg (status)::broker_translation_status
    );