-- +migrate Up

CREATE TABLE orders
(
    order_id    SERIAL PRIMARY KEY,
    discount_id INT REFERENCES discounts (id)
);