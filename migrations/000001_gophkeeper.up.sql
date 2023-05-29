CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
                        user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        login VARCHAR(255),
                        password VARCHAR(255)
);

CREATE TABLE users_data (
                       record_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id VARCHAR(256),
                       record_type INTEGER,
                       metadata VARCHAR(256),
                       encoded_data VARCHAR(256)
);