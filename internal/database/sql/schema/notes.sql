CREATE TABLE IF NOT EXISTS notes (
  id TEXT NOT NULL PRIMARY KEY,
  author TEXT NOT NULL,
  message TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  user_id TEXT NOT NULL DEFAULT "mmq9tll8as33g64",
  verified INTEGER NOT NULL DEFAULT 0
);

CREATE TRIGGER IF NOT EXISTS update_note_timestamp
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
  UPDATE notes
  SET updated_at = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
  WHERE id = NEW.id;
END;