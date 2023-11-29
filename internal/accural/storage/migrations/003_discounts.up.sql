-- +migrate Up

CREATE TABLE IF NOT EXISTS discounts
(
    id          SERIAL PRIMARY KEY,
    match       TEXT,
    reward      NUMERIC,
    reward_type TEXT
);