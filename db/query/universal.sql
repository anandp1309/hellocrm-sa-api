-- name: CreateUniversal :one
INSERT INTO mst_universal (
    universal_uuid,
    master_type_uuid,
    value_name,
    display_order,
    is_system,
    remarks,
    created_at,
    created_by_user_uuid,
    updated_at,
    updated_by_user_uuid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $7, $8
) RETURNING *;

-- name: GetUniversalsByMasterType :many
SELECT * FROM mst_universal
WHERE master_type_uuid = $1
ORDER BY display_order ASC;

-- name: UpdateUniversal :exec
UPDATE mst_universal SET value_name = $1, remarks = $2, updated_at = NOW() WHERE universal_uuid = $3;

-- name: DeactivateUniversal :exec
UPDATE mst_universal SET is_deleted = true, updated_at = NOW() WHERE universal_uuid = $1;
