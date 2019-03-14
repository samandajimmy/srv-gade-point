-- Table: users
-- for status 0 --> INACTIVE and 1 --> ACTIVE
-- for role 0 --> ADMIN and 1 --> USER
create table users (
  id SERIAL PRIMARY KEY NOT NULL,
  username VARCHAR(6) UNIQUE,
  email VARCHAR(255) UNIQUE,
  password  VARCHAR(255),
  status SMALLINT DEFAULT 0,
  role SMALLINT DEFAULT 0,
  updated_at TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NULL
);
CREATE INDEX index_users ON users (username, status, email);