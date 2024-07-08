CREATE TABLE users (
  id UUID PRIMARY KEY,
  firstname varchar NOT NULL,
  lastname varchar NOT NULL,
  email varchar UNIQUE NOT NULL,
  age smallint NOT NULL,
  created timestamptz NOT NULL
);