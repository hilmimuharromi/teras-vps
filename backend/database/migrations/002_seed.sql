-- TerasVPS Seed Data
-- Default pricing plans

-- Insert default plans
INSERT INTO plans (name, cores, memory, disk, price_monthly, price_daily, is_active) VALUES
('Starter', 1, 1024, 20, 50000, 1700, true),
('Standard', 2, 2048, 40, 100000, 3400, true),
('Premium', 4, 4096, 80, 200000, 6700, true)
ON CONFLICT (name) DO NOTHING;

-- Insert default admin user (password: admin123)
-- Hash generated with bcrypt: $2a$10$YourHashedPasswordHere
-- Note: This is a placeholder, use proper bcrypt hash in production
INSERT INTO users (username, email, password_hash, role, is_active) VALUES
('admin', 'admin@terasvps.com', '$2a$10$rZz7QZ8xXZxXZxXZxXZxZe', 'admin', true)
ON CONFLICT (email) DO NOTHING;
