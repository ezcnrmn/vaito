ALTER TABLE users ADD CONSTRAINT users_name_check CHECK (name ~ '^[a-zA-Z0-9]{3,50}$');

ALTER TABLE users ADD CONSTRAINT users_email_check CHECK (email ~ $$^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$$);


ALTER TABLE listings ADD CONSTRAINT listings_price_check CHECK (price >= 0);

ALTER TABLE listings ADD CONSTRAINT listings_title_check CHECK (LENGTH(title) BETWEEN 10 AND 300);

ALTER TABLE listings ADD CONSTRAINT listings_created_at_check CHECK (created_at < NOW() + INTERVAL '1 minute');

ALTER TABLE listings ADD CONSTRAINT listings_published_at_check CHECK (published_at >= created_at);
