CREATE TABLE devices (
    id INTEGER PRIMARY KEY,
    device_id TEXT,
    last_ping TEXT,
    is_admin INTEGER
);

CREATE TABLE notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    header TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'pending', 
    created_at TEXT DEFAULT (datetime('now'))
); --idk if this works for what you're tryna have but if it doesnt lmk!