-- Drop the index and column
DROP INDEX IF EXISTS idx_drawings_share_token;
ALTER TABLE drawings DROP COLUMN IF EXISTS share_token;
