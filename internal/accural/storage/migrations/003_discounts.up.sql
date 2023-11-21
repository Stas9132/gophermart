-- +migrate Up

CREATE TABLE discounts
(
    id          SERIAL PRIMARY KEY,
    match       TEXT,
    reward      NUMERIC,
    reward_type TEXT
);