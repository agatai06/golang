CREATE INDEX IF NOT EXISTS drones_title_idx ON drones USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS drones_categories_idx ON drones USING GIN (categories);