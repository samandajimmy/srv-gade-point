CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY NOT NULL,
    role_id INTEGER,
    user_image_path VARCHAR (300),
    first_name VARCHAR (300),
    second_name VARCHAR (300),
    email VARCHAR (300) UNIQUE,
    password VARCHAR (500),
    phone_number VARCHAR (300),
    address VARCHAR (300),
    updated_at BIGINT
);