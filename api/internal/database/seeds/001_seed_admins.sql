-- ============================================
-- Admin Users Seeder
-- Passwords are bcrypt hashed with cost 10
-- Cost 10 = 2^10 = 1024 hashing rounds (good balance of security and performance)
-- ============================================

-- Super Admin
-- Email: dev@susano.id
-- Password: admin1234
INSERT INTO admins (id, email, password, name, role, is_active, email_verified_at, created_at, updated_at)
VALUES (
    gen_uuid_v7(),
    'dev@susano.id',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'Yefta A. Wibowo',
    'super_admin',
    true,
    NOW(),
    NOW(),
    NOW()
);

-- Regular Admin
-- Email: admin@susano.id
-- Password: admin1234
INSERT INTO admins (id, email, password, name, role, is_active, email_verified_at, created_at, updated_at)
VALUES (
    gen_uuid_v7(),
    'admin@susano.id',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'Santi',
    'admin',
    true,
    NOW(),
    NOW(),
    NOW()
);

-- Cashier
-- Email: cashier@susano.id
-- Password: admin1234
INSERT INTO admins (id, email, password, name, role, is_active, email_verified_at, created_at, updated_at)
VALUES (
    gen_uuid_v7(),
    'cashier@susano.id',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'John Lennon',
    'cashier',
    true,
    NOW(),
    NOW(),
    NOW()
);
