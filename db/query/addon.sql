-- name: CreateAddon :one
INSERT INTO mst_addon (
    addon_uuid, name, category, addon_limit, price, status, icon_svg, icon_color
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetAddon :one
SELECT * FROM mst_addon
WHERE addon_uuid = $1 AND is_deleted = false;

-- name: ListAddons :many
SELECT * FROM mst_addon
WHERE is_deleted = false
ORDER BY created_at DESC;

-- name: UpdateAddon :one
UPDATE mst_addon
SET
    name = coalesce(sqlc.narg('name'), name),
    category = coalesce(sqlc.narg('category'), category),
    addon_limit = coalesce(sqlc.narg('addon_limit'), addon_limit),
    price = coalesce(sqlc.narg('price'), price),
    status = coalesce(sqlc.narg('status'), status),
    icon_svg = coalesce(sqlc.narg('icon_svg'), icon_svg),
    icon_color = coalesce(sqlc.narg('icon_color'), icon_color),
    updated_at = now()
WHERE addon_uuid = sqlc.arg('addon_uuid') AND is_deleted = false
RETURNING *;

-- name: DeleteAddon :exec
UPDATE mst_addon
SET is_deleted = true, updated_at = now()
WHERE addon_uuid = $1;
