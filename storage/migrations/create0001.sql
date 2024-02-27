-- DROP TABLE IF EXISTS Tags;
-- DROP TABLE IF EXISTS Events;

CREATE TABLE IF NOT EXISTS Tags (
    ID INTEGER PRIMARY KEY,
    Tag TEXT UNIQUE NOT NULL,
    Comment TEXT
);

CREATE TABLE IF NOT EXISTS Events (
    ID INTEGER PRIMARY KEY,
    Tag INTEGER NOT NULL,
    Date INTEGER NOt NULL,
    Comment TEXT,

    FOREIGN KEY (Tag)
        REFERENCES Tags(ID)
);