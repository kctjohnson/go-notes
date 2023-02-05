DROP TABLE IF EXISTS tags;

CREATE TEMPORARY TABLE temp AS
SELECT id, title, content, created_date, last_edited_date FROM notes;

DROP TABLE notes;

CREATE TABLE notes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TINYTEXT NOT NULL,
  content MEDIUMTEXT NOT NULL,
  created_date DATETIME NOT NULL,
  last_edited_date DATETIME NOT NULL
);

INSERT INTO notes (id, title, content, created_date, last_edited_date)
SELECT id, title, content, created_date, last_edited_date FROM temp;

DROP TABLE temp;
