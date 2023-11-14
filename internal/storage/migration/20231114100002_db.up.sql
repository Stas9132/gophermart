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