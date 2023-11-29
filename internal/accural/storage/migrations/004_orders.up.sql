-- +migrate Up

CREATE TABLE discounts
(
    id          SERIAL PRIMARY KEY,
    match       TEXT,
    reward      NUMERIC,
    reward_type TEXT
);

CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    order_id    VARCHAR(255),
    discount_id INT REFERENCES discounts (id)
);