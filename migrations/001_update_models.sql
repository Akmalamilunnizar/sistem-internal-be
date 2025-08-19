-- Migration: Update models with new fields
-- This migration adds the new fields to existing tables

-- Update customers table
ALTER TABLE customers 
ADD COLUMN IF NOT EXISTS address TEXT,
ADD COLUMN IF NOT EXISTS gps_lat DECIMAL(10,8),
ADD COLUMN IF NOT EXISTS gps_long DECIMAL(11,8);

-- Update employees table
ALTER TABLE employees 
ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active';

-- Update trouble_tickets table
ALTER TABLE trouble_tickets 
ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'connection_issue',
MODIFY COLUMN status VARCHAR(20) NOT NULL DEFAULT 'open';

-- Update status values for existing tickets
UPDATE trouble_tickets SET status = 'open' WHERE status = 'received';

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees(status);
CREATE INDEX IF NOT EXISTS idx_tickets_customer_id ON trouble_tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON trouble_tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_type ON trouble_tickets(type);
