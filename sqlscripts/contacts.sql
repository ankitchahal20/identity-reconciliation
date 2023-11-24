CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    phoneNumber TEXT,
    email TEXT,
    linkedId INTEGER,
    linkPrecedence TEXT,
    createdAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMPTZ
);