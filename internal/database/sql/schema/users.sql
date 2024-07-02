CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  verified INTEGER NOT NULL DEFAULT 0
);

CREATE TRIGGER IF NOT EXISTS update_user_timestamp
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
  UPDATE users
  SET updated_at = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
  WHERE id = NEW.id;
END;