CREATE TABLE IF NOT EXISTS swift_codes (
    id SERIAL PRIMARY KEY,
    swift_code VARCHAR(11) UNIQUE NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    country_iso2 CHAR(2) NOT NULL,
    country_name VARCHAR(255) NOT NULL,
    is_headquarter BOOLEAN NOT NULL
);
