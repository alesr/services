CREATE TABLE IF NOT EXISTS email_verifications (
    token VARCHAR(255) NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX ON email_verifications(token);
