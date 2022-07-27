CREATE TABLE link_visits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    id_link INTEGER,
    ts TEXT,
    FOREIGN KEY(id_link) REFERENCES links(id)
);