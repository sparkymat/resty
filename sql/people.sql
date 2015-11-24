-- name: create-people-table
CREATE TABLE people (
  id  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name VARCHAR(255),
  email VARCHAR(255),
);

-- name: create-person
INSERT INTO people (name,  email) VALUES (?, ?)

-- name: find-person-by-id
SELECT * FROM people WHERE id = ? LIMIT 1

-- name: drop-people-table
DROP TABLE people
