-- name: GetAllVerifiedTranslations :many
SELECT
    id,
    source_text,
    target_text
FROM
    broker_translations
WHERE
    status = 'verified';

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