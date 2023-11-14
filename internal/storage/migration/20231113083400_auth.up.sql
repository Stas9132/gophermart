CREATE TABLE IF NOT EXISTS auth (
    id serial PRIMARY KEY,
    login varchar(255) unique not null ,
    password varchar(255));
