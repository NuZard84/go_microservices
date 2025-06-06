CREATE TABLE IF NOT EXISTS orders (
    id CHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    account_id CHAR(27) NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL,
);

CREATE TABLE IF NOT EXISTS order_products (
    order_id CHAR(27) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id CHAR(27) NOT NULL,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (order_id, product_id)
);

