CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE company_type AS ENUM (
    'Corporations',
    'NonProfit',
    'Cooperative',
    'Sole Proprietorship');