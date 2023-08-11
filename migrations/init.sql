-- migrations/init.sql
CREATE TABLE articles (
  id SERIAL PRIMARY KEY,
  author TEXT,
  title TEXT,
  name TEXT,
  body TEXT,
  created TIMESTAMP
);
