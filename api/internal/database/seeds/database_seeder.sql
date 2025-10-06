-- ============================================
-- Database Seeder
-- This file orchestrates all seeders
-- ============================================

\echo '================================'
\echo 'Starting database seeding...'
\echo '================================'

\echo ''
\echo '-> Seeding admins...'
\i internal/database/seeds/001_seed_admins.sql

\echo ''
\echo '================================'
\echo 'Seeding completed successfully!'
\echo '================================'
