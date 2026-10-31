CREATE TABLE users (
  id            BIGINT GENERATED ALWAYS AS IDENTITY
                PRIMARY KEY,
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL
);

CREATE TABLE books (
  id       BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title    TEXT NOT NULL,
  author   TEXT NOT NULL,
  year     INTEGER NOT NULL,
  added_by BIGINT NOT NULL REFERENCES users(id)
);

CREATE TABLE loans (
  id          BIGINT GENERATED ALWAYS AS IDENTITY
              PRIMARY KEY,
  book_id     BIGINT NOT NULL
              REFERENCES books(id) ON DELETE CASCADE,
  user_id     BIGINT NOT NULL REFERENCES users(id),
  borrowed_at TIMESTAMPTZ NOT NULL,
  due_at      TIMESTAMPTZ NOT NULL,
  returned_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX one_open_loan
  ON loans(book_id) WHERE returned_at IS NULL;
