ALTER TABLE products ADD CONSTRAINT products_price_check CHECK (price > 0);
ALTER TABLE products ADD CONSTRAINT products_year_check CHECK (year BETWEEN 2000 AND date_part('year', now()));
ALTER TABLE products ADD CONSTRAINT categories_length_check CHECK (array_length(categories, 1) BETWEEN 1 AND 5);