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
CREATE TABLE IF NOT EXISTS orders (
    id serial PRIMARY KEY,
    number varchar(255) unique not null ,
    status varchar(255) ,
    accrual  INTEGER,
    uploaded_at timestamp);