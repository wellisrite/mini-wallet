CREATE DATABASE mini_wallet;

CREATE TABLE wallet (
    id UUID PRIMARY KEY,
    owned_by UUID,
    status VARCHAR(255),
    enabled_at TIMESTAMP,
    balance INT,
    token VARCHAR(255)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,          -- Unique identifier for the transaction
    customer_id UUID NOT NULL,    -- ID of the customer associated with the transaction
    amount INT NOT NULL,         -- Amount of the transaction
    transacted_at TIMESTAMP NOT NULL, -- Timestamp of the transaction
    status VARCHAR(255) NOT NULL, -- Status of the transaction (e.g., 'success', 'failed')
    type VARCHAR(255) NOT NULL,   -- Type of the transaction (e.g., 'deposit', 'withdrawal')
    reference_id UUID NOT NULL    -- Unique reference ID for the transaction
);

ALTER TABLE transactions
ADD CONSTRAINT unique_reference_id UNIQUE (reference_id);