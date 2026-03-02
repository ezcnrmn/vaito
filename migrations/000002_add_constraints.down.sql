ALTER TABLE users DROP CONSTRAINT IF EXISTS users_name_check;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_check;

ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_price_check;
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_title_check;
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_created_at_check;
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_published_at_check;


