CREATE SCHEMA app;

CREATE TABLE app.users (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    login VARCHAR(100) NOT NULL CHECK(char_length(login) BETWEEN 8 AND 100),
    password VARCHAR(100) NOT NULL CHECK(char_length(password) BETWEEN 8 AND 100),
    city VARCHAR(100) NOT NULL CHECK(
        char_length(city) BETWEEN 1 AND 100
        AND
        city ~ '^[A-Z]'
    ),

    UNIQUE(login)
);

CREATE TABLE app.couriers (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    login VARCHAR(100) NOT NULL CHECK(char_length(login) BETWEEN 8 AND 100),
    password VARCHAR(100) NOT NULL CHECK(char_length(password) BETWEEN 8 AND 100),
    city VARCHAR(100) NOT NULL CHECK(
        char_length(city) BETWEEN 1 AND 100
        AND
        city ~ '^[A-Z]'
    ),
    orders_complete BIGINT NOT NULL,
    is_free BOOLEAN NOT NULL,

    UNIQUE(login)
);

CREATE TABLE app.orders (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(20) NOT NULL CHECK(char_length(name) BETWEEN 3 AND 20),
    price BIGINT NOT NULL,
    is_complete BOOLEAN NOT NULL,
    user_id INTEGER NOT NULL,
    courier_id INTEGER,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES app.users(id),
    CONSTRAINT fk_courier_id FOREIGN KEY (courier_id) REFERENCES app.couriers(id)
);