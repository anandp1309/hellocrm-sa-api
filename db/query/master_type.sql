-- name: CreateMasterType :one
INSERT INTO mst_master_type (
    master_type_uuid,
    master_type_name,
    display_order,
    is_system,
    remarks,
    created_at,
    created_by_user_uuid,
    updated_at,
    updated_by_user_uuid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $6, $7
) RETURNING *;

-- name: GetMasterType :one
SELECT * FROM mst_master_type
WHERE master_type_uuid = $1 AND is_deleted = false;

-- name: ListMasterTypes :many
SELECT * FROM mst_master_type
WHERE is_deleted = false
ORDER BY display_order ASC
LIMIT $1 OFFSET $2;

-- name: UpdateMasterType :one
UPDATE mst_master_type
SET
    master_type_name = COALESCE(sqlc.narg('master_type_name'), master_type_name),
    display_order = COALESCE(sqlc.narg('display_order'), display_order),
    is_system = COALESCE(sqlc.narg('is_system'), is_system),
    remarks = COALESCE(sqlc.narg('remarks'), remarks),
    updated_at = NOW(),
    updated_by_user_uuid = $2
WHERE master_type_uuid = $1 AND is_deleted = false
RETURNING *;

-- name: DeleteMasterType :exec
UPDATE mst_master_type
SET
    is_deleted = true,
    deleted_at = NOW(),
    deleted_by_user_uuid = $2
WHERE master_type_uuid = $1 AND is_deleted = false;
