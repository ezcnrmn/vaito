ALTER TABLE users ADD CONSTRAINT users_name_check CHECK (name ~ '^[a-zA-Z0-9]{3,50}$');

ALTER TABLE users ADD CONSTRAINT users_email_check CHECK (email ~ $$^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$$);
