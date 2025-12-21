-- Add share_token column to drawings table for public sharing
ALTER TABLE drawings ADD COLUMN share_token VARCHAR(255) NULL;

-- Create unique index on share_token for fast lookups and uniqueness constraint
CREATE UNIQUE INDEX idx_drawings_share_token ON drawings(share_token) WHERE share_token IS NOT NULL;
