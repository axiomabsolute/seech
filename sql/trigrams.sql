-- name: create-table
CREATE TABLE IF NOT EXISTS sources (
    id INTEGER PRIMARY KEY,
    file_path NOT NULL,
    line_number NOT NULL,
    doc NOT NULL,
    UNIQUE (file_path, line_number)
);

-- name: create-trigram-table
CREATE VIRTUAL TABLE trigrams USING fts5(
    doc, content='sources', content_rowid='id', tokenize="trigram"
);

-- name: create-insert-trigger
CREATE TRIGGER sources_ai AFTER INSERT ON sources BEGIN
    INSERT INTO trigrams(rowid, doc) VALUES (new.id, new.doc);
END;

-- name: create-delete-trigger
CREATE TRIGGER sources_ad AFTER DELETE ON sources BEGIN
    INSERT INTO trigrams(trigrams, rowid, doc) VALUES('delete', old.id, old.doc);
END;

-- name: create-update-trigger
CREATE TRIGGER sources_au AFTER UPDATE ON sources BEGIN
    INSERT INTO trigrams(trigrams, rowid, doc) VALUES('delete', old.id, old.doc);
    INSERT INTO trigrams(rowid, doc) VALUES (new.id, new.doc);
END;
