CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target TEXT,
    redirect INTEGER DEFAULT 302
);