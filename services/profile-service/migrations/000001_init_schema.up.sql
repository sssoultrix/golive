CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    bio TEXT,
    image TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_profiles_login ON profiles (login);
CREATE INDEX idx_profiles_email ON profiles (email);

