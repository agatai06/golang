ALTER TABLE drones ADD CONSTRAINT drones_price_check CHECK (price > 0);
ALTER TABLE drones ADD CONSTRAINT drones_year_check CHECK (year BETWEEN 2000 AND date_part('year', now()));
ALTER TABLE drones ADD CONSTRAINT materials_length_check CHECK (array_length(categories, 1) BETWEEN 1 AND 5);