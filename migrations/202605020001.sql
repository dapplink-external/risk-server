DO
$$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'uint256') THEN
        CREATE DOMAIN UINT256 AS NUMERIC
            CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
    ELSE
        ALTER DOMAIN UINT256 DROP CONSTRAINT uint256_check;
        ALTER DOMAIN UINT256 ADD
        CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS addresses (
    guid         VARCHAR PRIMARY KEY,
    address      VARCHAR UNIQUE NOT NULL,
    address_type VARCHAR(10) NOT NULL DEFAULT 'white',  -- black:黑地址; gray:灰地址
    created_at   TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT   check_address_type CHECK (address_type IN ('white', 'black', 'gray'))
);
CREATE INDEX IF NOT EXISTS idx_addresses_address ON addresses (address);

CREATE TABLE IF NOT EXISTS transactions (
    guid          VARCHAR PRIMARY KEY,
    request_id    VARCHAR NOT NULL,
    from_address  VARCHAR NOT NULL,
    to_address    VARCHAR NOT NULL,
    token_address VARCHAR NOT NULL,
    token_id      VARCHAR NOT NULL,
    token_meta    VARCHAR NOT NULL,
    fee           UINT256 NOT NULL,
    amount        UINT256 NOT NULL,
    status        VARCHAR NOT NULL,
    tx_type       VARCHAR NOT NULL,
    created_at    TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT   check_status CHECK (status IN ('pending', 'checked_fail', 'checked_pass'))
);
CREATE INDEX IF NOT EXISTS transactions_request_id ON transactions (request_id);
