CREATE TABLE IF NOT EXISTS User (
    Username            TEXT PRIMARY KEY,
    firstName           TEXT NOT NULL,
    lastName            TEXT NOT NULL,
    gender              INTEGER NOT NULL CHECK (gender IN (0, 1)), -- 1 for female, 0 for male
    age                 INTEGER NOT NULL,
    email               TEXT NOT NULL UNIQUE,
    password            TEXT NOT NULL,
    sessionToken        TEXT,
    sessionExpiration   DATETIME
);

-- CREATE TABLE IF NOT EXISTS Category (
--     CategoryID      INTEGER PRIMARY KEY AUTOINCREMENT,
--     CategoryName    TEXT NOT NULL UNIQUE
-- );

CREATE TABLE IF NOT EXISTS Post (
    PostID          INTEGER PRIMARY KEY AUTOINCREMENT,
    Title           TEXT NOT NULL,
    Content         TEXT NOT NULL,
    Category        TEXT NOT NULL,
    CreatedDate     DATETIME DEFAULT CURRENT_TIMESTAMP,
    username        TEXT,
    FOREIGN KEY (username) REFERENCES User(username)
);

CREATE TABLE IF NOT EXISTS Comment (
    CommentID       INTEGER PRIMARY KEY AUTOINCREMENT,
    Content         TEXT NOT NULL,
    CreatedDate     DATETIME DEFAULT CURRENT_TIMESTAMP,
    PostID          INTEGER,
    username        TEXT,
    FOREIGN KEY (PostID) REFERENCES Post(PostID),
    FOREIGN KEY (username) REFERENCES User(username)
);

-- CREATE TABLE IF NOT EXISTS PostCategory (
--     PostCategoryID  INTEGER PRIMARY KEY AUTOINCREMENT,
--     PostID          INTEGER,
--     CategoryID      INTEGER,
--     FOREIGN KEY (PostID) REFERENCES Post(PostID),
--     FOREIGN KEY (CategoryID) REFERENCES Category(CategoryID)
-- );

CREATE TABLE IF NOT EXISTS Interaction (
    InteractionID   INTEGER PRIMARY KEY AUTOINCREMENT,
    PostID          INTEGER,
    username        TEXT,
    Kind            INTEGER NOT NULL CHECK (Kind IN (0, 1)), -- 1 for like, 0 for dislike
    FOREIGN KEY (PostID) REFERENCES Post(PostID),
    FOREIGN KEY (username) REFERENCES User(username),
    UNIQUE (PostID, username) 
);

CREATE TABLE IF NOT EXISTS CommentInteractions (
    CommentInteractionID    INTEGER PRIMARY KEY AUTOINCREMENT,
    CommentID               INTEGER,
    username                TEXT,
    Kind                    INTEGER NOT NULL CHECK (Kind IN (0, 1)), -- 1 for like, 0 for dislike
    FOREIGN KEY (CommentID) REFERENCES Comment(CommentID),
    FOREIGN KEY (username) REFERENCES User(username),
    UNIQUE (CommentID, username) 
);

CREATE TABLE IF NOT EXISTS Chat (
    ChatID          INTEGER PRIMARY KEY AUTOINCREMENT,
    SenderUsername     TEXT,
    RecipientUsername  TEXT,
    Content            TEXT NOT NULL,
    CreatedDate        DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (SenderUsername) REFERENCES User(username),
    FOREIGN KEY (RecipientUsername) REFERENCES User(username)
);

-- THIS WAS CAUSING AN ERROR BECAUSE IT WAS EMPTY!!

-- CREATE TABLE IF NOT EXISTS Notifications (
--     -- add later
-- );