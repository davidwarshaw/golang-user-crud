-- Create the schema on container startup
-- (NOTE: we are the superuser in the default DB)

DROP TABLE IF EXISTS user_accounts CASCADE;
CREATE TABLE user_accounts (
    id SERIAL PRIMARY KEY,
    
    user_name VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    first_name VARCHAR(1024),
    middle_name VARCHAR(1024),
    last_name VARCHAR(1024),
    email VARCHAR(1024),
    primary_phone_number VARCHAR(17),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
