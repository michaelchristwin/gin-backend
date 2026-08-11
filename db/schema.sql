-- ============================================================
-- Sync state
-- ============================================================

CREATE TABLE sync_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),

    -- Last block that was successfully persisted.
    last_synced_block INTEGER NOT NULL DEFAULT 0,

    last_synced_at INTEGER,

    -- Useful for monitoring the health of the sync process.
    last_error TEXT
);

INSERT INTO sync_state (id, last_synced_block)
VALUES (1, 0);


-- ============================================================
-- Raw meter data
-- ============================================================

CREATE TABLE meter_data_points (
    transaction_id TEXT PRIMARY KEY,

    meter_number INTEGER NOT NULL,

    -- Unix timestamp.
    timestamp INTEGER NOT NULL,

    nonce INTEGER,

    voltage REAL,
    energy REAL,

    longitude REAL,
    latitude REAL,

    signature TEXT,
    public_key TEXT,

    -- The block this datapoint came from.
    block_number INTEGER NOT NULL,

    created_at INTEGER NOT NULL DEFAULT (unixepoch())
);


-- ============================================================
-- Raw data indexes
-- ============================================================

CREATE INDEX idx_meter_data_meter_timestamp
ON meter_data_points (meter_number, timestamp);

CREATE INDEX idx_meter_data_timestamp
ON meter_data_points (timestamp);

CREATE INDEX idx_meter_data_block
ON meter_data_points (block_number);


-- ============================================================
-- Hourly aggregates
-- ============================================================

CREATE TABLE hourly_meter_stats (
    meter_number INTEGER NOT NULL,

    -- Start of the hour as Unix timestamp.
    hour_start INTEGER NOT NULL,

    -- Number of datapoints contributing to this bucket.
    sample_count INTEGER NOT NULL DEFAULT 0,

    energy_sum REAL NOT NULL DEFAULT 0,

    voltage_sum REAL NOT NULL DEFAULT 0,

    voltage_min REAL,
    voltage_max REAL,

    PRIMARY KEY (meter_number, hour_start)
);


-- ============================================================
-- Daily aggregates
-- ============================================================

CREATE TABLE daily_meter_stats (
    meter_number INTEGER NOT NULL,

    -- Start of day as Unix timestamp.
    day_start INTEGER NOT NULL,

    sample_count INTEGER NOT NULL DEFAULT 0,

    energy_sum REAL NOT NULL DEFAULT 0,

    voltage_sum REAL NOT NULL DEFAULT 0,

    voltage_min REAL,
    voltage_max REAL,

    PRIMARY KEY (meter_number, day_start)
);


-- ============================================================
-- Monthly aggregates
-- ============================================================

CREATE TABLE monthly_meter_stats (
    meter_number INTEGER NOT NULL,

    -- Start of month as Unix timestamp.
    month_start INTEGER NOT NULL,

    sample_count INTEGER NOT NULL DEFAULT 0,

    energy_sum REAL NOT NULL DEFAULT 0,

    voltage_sum REAL NOT NULL DEFAULT 0,

    voltage_min REAL,
    voltage_max REAL,

    PRIMARY KEY (meter_number, month_start)
);


-- ============================================================
-- Yearly aggregates
-- ============================================================

CREATE TABLE yearly_meter_stats (
    meter_number INTEGER NOT NULL,

    -- Start of year as Unix timestamp.
    year_start INTEGER NOT NULL,

    sample_count INTEGER NOT NULL DEFAULT 0,

    energy_sum REAL NOT NULL DEFAULT 0,

    voltage_sum REAL NOT NULL DEFAULT 0,

    voltage_min REAL,
    voltage_max REAL,

    PRIMARY KEY (meter_number, year_start)
);


-- ============================================================
-- Aggregate indexes
-- ============================================================

CREATE INDEX idx_hourly_meter_hour
ON hourly_meter_stats (meter_number, hour_start);

CREATE INDEX idx_daily_meter_day
ON daily_meter_stats (meter_number, day_start);

CREATE INDEX idx_monthly_meter_month
ON monthly_meter_stats (meter_number, month_start);

CREATE INDEX idx_yearly_meter_year
ON yearly_meter_stats (meter_number, year_start);


-- ============================================================
-- SQLite configuration
-- ============================================================

PRAGMA journal_mode = WAL;

PRAGMA synchronous = NORMAL;

PRAGMA foreign_keys = ON;