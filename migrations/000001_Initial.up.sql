CREATE TABLE IF NOT EXISTS notes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title TINYTEXT NOT NULL,
  content MEDIUMTEXT NOT NULL,
  created_date DATETIME NOT NULL,
  last_edited_date DATETIME NOT NULL
);
