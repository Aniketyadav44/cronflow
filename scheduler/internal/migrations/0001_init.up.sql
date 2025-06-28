CREATE TABLE IF NOT EXISTS jobs(
    id SERIAL PRIMARY KEY,
    cron_id INT,
    cron_expr TEXT NOT NULL,
    -- ping, email, slack, webhook
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_entries(
    id SERIAL PRIMARY KEY,
    job_id INT,
    -- running, completed, failed, permanently_failed
    status TEXT,
    retries INT DEFAULT 0,
    output TEXT,
    error TEXT,
    scheduled_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
