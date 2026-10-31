CREATE TABLE IF NOT EXISTS users (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS books (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  title    TEXT NOT NULL,
  author   TEXT NOT NULL,
  year     INTEGER NOT NULL,
  added_by INTEGER NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS loans (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  book_id     INTEGER NOT NULL
              REFERENCES books(id) ON DELETE CASCADE,
  user_id     INTEGER NOT NULL REFERENCES users(id),
  borrowed_at DATETIME NOT NULL,
  due_at      DATETIME NOT NULL,
  returned_at DATETIME
);

CREATE UNIQUE INDEX IF NOT EXISTS one_open_loan
  ON loans(book_id) WHERE returned_at IS NULL;
