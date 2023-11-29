CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(50),
                                     email VARCHAR(50),
                                     password VARCHAR(50)
);
CREATE TABLE IF NOT EXISTS auth (
                                    id serial PRIMARY KEY,
                                    login varchar(255) unique not null ,
                                    password varchar(255));
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