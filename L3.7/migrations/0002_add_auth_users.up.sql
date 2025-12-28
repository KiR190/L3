-- Users table for JWT auth (email/password) + roles
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin','manager','viewer')) DEFAULT 'viewer',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- admin / Admin12345
INSERT INTO users (email, password_hash, role)
VALUES
('admin@local', '$2a$10$Qlf1rKroBijkouRjTu16d.7KPdcLK8p6rLbEcXC1x0.dbO2SLyMuO', 'admin')
ON CONFLICT (email) DO NOTHING;

-- bob(manager) / Manager12345
INSERT INTO users (email, password_hash, role)
VALUES
('bob@local', '$2a$10$fzP3aWb0ZvZpvr02rHtw/eieHGGJjD.UTWdcNrH2og0TAXMZKBh2G', 'manager')
ON CONFLICT (email) DO NOTHING;

-- alice(viewer) / Viewer12345
INSERT INTO users (email, password_hash, role)
VALUES
('alice@local', '$2a$10$hkLVWHvi2PPWD.yAdbn8auT1K3WHQdGKMvO./WDFXeAuZy21xR3.e', 'viewer')
ON CONFLICT (email) DO NOTHING;


