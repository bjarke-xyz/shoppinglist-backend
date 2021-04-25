CREATE TABLE items (
  id INT PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  name VARCHAR(255) NOT NULL,
  owner_id VARCHAR(255) NOT NULL,
)
CREATE INDEX active_items ON items(deleted_at) where deleted_at IS NOT NULL
