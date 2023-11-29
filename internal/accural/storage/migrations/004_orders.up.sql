-- +migrate Up

CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    order_id    VARCHAR(255),
    discount_id INT REFERENCES discounts (id)
);