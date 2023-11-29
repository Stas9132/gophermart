-- +migrate Up

CREATE TABLE IF NOT EXISTS discounts
(
    id          SERIAL PRIMARY KEY,
    match       TEXT,
    reward      NUMERIC,
    reward_type TEXT
);

CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    order_id    VARCHAR(255),
    discount_id INT
);