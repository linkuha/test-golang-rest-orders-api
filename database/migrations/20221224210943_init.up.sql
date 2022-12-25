CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    username varchar NOT NULL UNIQUE,
    password_hash varchar NOT NULL
--     roles jsonb NOT NULL DEFAULT '{"user"}'
--     status smallint NOT NULL,
--     created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
--     updated_at TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT NULL,
);

CREATE TABLE IF NOT EXISTS user_profiles (
    id serial NOT NULL UNIQUE,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    middle_name varchar,
    sex char(1) NOT NULL,
    age smallint
);

CREATE TABLE IF NOT EXISTS products (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name varchar NOT NULL,
    description varchar,
    left_in_stock int NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS product_prices (
    id bigserial NOT NULL UNIQUE,
    product_id uuid REFERENCES products(id) ON DELETE CASCADE NOT NULL,
    currency char(3) NOT NULL,
    price money NOT NULL DEFAULT 0
--     status smallint NOT NULL,
);
CREATE UNIQUE INDEX uq_product_currency on product_prices (product_id, currency);

CREATE TABLE IF NOT EXISTS user_orders (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    number bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS user_order_products (
    id bigserial NOT NULL UNIQUE,
    order_id uuid NOT NULL,
    product_id uuid NOT NULL,
    amount int NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_order_product on user_order_products (order_id, product_id);
ALTER TABLE user_order_products ADD CONSTRAINT FK_userOrderProducts_orderId
    FOREIGN KEY (order_id) REFERENCES user_orders(id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE user_order_products ADD CONSTRAINT FK_userOrderProducts_productId
    FOREIGN KEY (product_id) REFERENCES products(id) NOT DEFERRABLE INITIALLY IMMEDIATE;

CREATE TABLE IF NOT EXISTS user_followers (
    id bigserial NOT NULL UNIQUE,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    follower_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
--     status smallint NOT NULL,
    PRIMARY KEY (user_id, follower_id)
);
