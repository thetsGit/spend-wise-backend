-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    session_token VARCHAR(255) UNIQUE,
    expires_at TIMESTAMPTZ,
    oauth_id VARCHAR(255) UNIQUE NOT NULL,
    oauth_email VARCHAR(255) UNIQUE NOT NULL,
    oauth_name VARCHAR(255) UNIQUE NOT NULL,
    oauth_picture TEXT NOT NULL,
    oauth_access_token TEXT,
    oauth_refresh_token TEXT,
    oauth_token_expiry TIMESTAMPTZ,
    oauth_token_type VARCHAR(255),
    oauth_scope TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create email table
CREATE TABLE email (
	id SERIAL PRIMARY KEY,
	sender VARCHAR(255),
	recipient VARCHAR(255),
	subject TEXT,
	body TEXT,
	date TIMESTAMPTZ,
	status VARCHAR(20) NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	UNIQUE (sender, recipient, subject, date)
);

-- Create spending table
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
    email_id INTEGER NOT NULL REFERENCES email(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Create saas_discovery table
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
    email_id INTEGER NOT NULL REFERENCES email(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);


-- Adding indexes - email relation

CREATE INDEX email_user_id ON email(user_id);
CREATE INDEX email_date ON email(date);
CREATE INDEX email_status ON email(status);

-- Adding indexes - spending relation

CREATE INDEX spending_user_id ON spending(user_id);
CREATE INDEX spending_email_id ON spending(email_id);
CREATE INDEX spending_category ON spending(category);
CREATE INDEX spending_transaction_date ON spending(transaction_date);

-- Adding indexes - saas_discovery relation

CREATE INDEX saas_discovery_user_id ON saas_discovery(user_id);
CREATE INDEX saas_discovery_email_id ON saas_discovery (email_id);
CREATE INDEX saas_discovery_product_name ON saas_discovery(product_name);
CREATE INDEX saas_discovery_signal_type ON saas_discovery(signal_type);
