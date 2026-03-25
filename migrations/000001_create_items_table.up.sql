CREATE TABLE IF NOT EXISTS items (
  id bigserial PRIMARY KEY,
  name text NOT NULL,
  rarity text NOT NULL,
  price integer NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
)