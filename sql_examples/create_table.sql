DROP TABLE IF EXISTS accounts;
CREATE TABLE IF NOT EXISTS accounts (
  user_id serial PRIMARY KEY,
  username VARCHAR ( 50 ) UNIQUE NOT NULL,
  password VARCHAR ( 50 ) NOT NULL,
  email VARCHAR ( 255 ) UNIQUE NOT NULL,
  first_name  VARCHAR ( 255 ) UNIQUE NOT NULL,
  middle_name VARCHAR ( 255 ) UNIQUE NOT NULL,
  last_name   VARCHAR ( 255 ) UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  last_login_at TIMESTAMP
);
