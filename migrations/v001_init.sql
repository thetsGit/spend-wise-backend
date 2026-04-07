
CREATE TABLE email (
    id SERIAL PRIMARY KEY,
    sender VARCHAR(255),
    recipient VARCHAR(255),
    subject TEXT,
    body TEXT,
    date TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (sender, recipient, subject, date)
);

CREATE TABLE spending (
    id SERIAL PRIMARY KEY,
    merchant VARCHAR(255),
    amount DECIMAL(12,2),
    currency VARCHAR(20),
    category VARCHAR(255),
    transaction_date TIMESTAMPTZ,
    ai_confidence VARCHAR(255),
    confidence VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    email_id INTEGER NOT NULL REFERENCES email(id) ON DELETE CASCADE
);

CREATE TABLE saas_discovery (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255),
    signal_type VARCHAR(255),
    billing_cycle VARCHAR(255),
    estimated_cost DECIMAL(12, 2),
    currency VARCHAR(20),
    ai_confidence VARCHAR(255),
    confidence VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    email_id INTEGER NOT NULL REFERENCES email(id) ON DELETE CASCADE
);

-- Adding indexes - email relation

CREATE INDEX email_date ON email(date);
CREATE INDEX email_status ON email(status);

-- Adding indexes - spending relation

CREATE INDEX spending_email_id ON spending(email_id);
CREATE INDEX spending_category ON spending(category);
CREATE INDEX spending_transaction_date ON spending(transaction_date);

-- Adding indexes - saas_discovery relation

CREATE INDEX saas_discovery_email_id ON saas_discovery(email_id);
CREATE INDEX saas_discovery_product_name ON saas_discovery(product_name);
CREATE INDEX saas_discovery_signal_type ON saas_discovery(signal_type);
