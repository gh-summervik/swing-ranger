DROP TABLE IF EXISTS public.eod_prices CASCADE;

CREATE TABLE IF NOT EXISTS public.eod_prices (
    symbol TEXT NOT NULL,
    date_eod DATE NOT NULL,
    open NUMERIC(22,4) NOT NULL,
    high NUMERIC(22,4) NOT NULL,
    low NUMERIC(22,4) NOT NULL,
    close NUMERIC(22,4) NOT NULL,
    volume DOUBLE PRECISION NOT NULL,
    created_by TEXT NOT NULL,
    updated_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    created_at_unix_ms BIGINT NOT NULL,
    updated_at_unix_ms BIGINT NOT NULL,
    PRIMARY KEY (symbol, date_eod)
);
