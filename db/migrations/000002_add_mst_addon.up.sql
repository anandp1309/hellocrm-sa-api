CREATE TABLE mst_addon (
    addon_uuid UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    addon_limit VARCHAR(100) NOT NULL,
    price VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Active',
    icon_svg TEXT DEFAULT '',
    icon_color VARCHAR(50) DEFAULT 'text-indigo-600',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    is_deleted BOOLEAN NOT NULL DEFAULT false
);
