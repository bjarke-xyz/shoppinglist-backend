CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  name VARCHAR(255) NOT NULL,
  owner_id VARCHAR(36) NOT NULL
)
