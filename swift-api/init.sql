CREATE TABLE IF NOT EXISTS swift_codes (
    id SERIAL PRIMARY KEY,
    swift_code TEXT UNIQUE NOT NULL,
    bank_name TEXT NOT NULL,
    address TEXT NOT NULL,
    country_iso2 TEXT NOT NULL,
    country_name TEXT NOT NULL,
    is_headquarter BOOLEAN NOT NULL
);
