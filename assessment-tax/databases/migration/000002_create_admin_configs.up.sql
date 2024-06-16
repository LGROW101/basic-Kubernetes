BEGIN;

CREATE TABLE
    admin_configs (
        id SERIAL PRIMARY KEY,
        personal_deduction DECIMAL(10, 2) NOT NULL DEFAULT 60000.00 CHECK (
            personal_deduction >= 0
            AND personal_deduction <= 100000
        ),
        k_receipt DECIMAL(10, 2) NOT NULL DEFAULT 50000.00 CHECK (
            k_receipt >= 0
            AND k_receipt <= 100000
        ),
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

-- Memasukkan data awal ke dalam tabel admin_configs
INSERT INTO
    admin_configs (personal_deduction, k_receipt)
VALUES
    (60000.00, 50000.00);

CREATE TABLE
    tax_calculations (
        id SERIAL PRIMARY KEY,
        totalIncome DECIMAL(10, 2) NOT NULL,
        wht DECIMAL(10, 2) NOT NULL DEFAULT '0.00',
        personal_allowance DECIMAL(10, 2) NOT NULL DEFAULT '60000.00',
        donation DECIMAL(10, 2) NOT NULL DEFAULT '0.00',
        k_receipt DECIMAL(10, 2) NOT NULL DEFAULT '0.00',
        tax DECIMAL(10, 2) NOT NULL DEFAULT '0.00',
        created_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_tax_calculations_created_at ON tax_calculations (created_at);

COMMIT;