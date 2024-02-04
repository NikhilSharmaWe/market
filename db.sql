CREATE TABLE products(
    product_name TEXT NOT NULL PRIMARY KEY,
	category TEXT NOT NULL,
	price INT NOT NULL
);

CREATE TABLE users(
    username TEXT NOT NULL PRIMARY KEY,
	password_hash TEXT NOT NULL,
	email TEXT NOT NULL
);

CREATE TABLE admins(
    username TEXT NOT NULL PRIMARY KEY REFERENCES users(username) ON DELETE CASCADE
);

CREATE TABLE orders(
    order_id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
	total_cost INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE order_per_product(
	order_id TEXT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
	product_name TEXT NOT NULL,
	quantity INT NOT NULL,
	PRIMARY KEY (order_id, product_name)
);

CREATE TABLE inventory(
    product_name TEXT NOT NULL PRIMARY KEY REFERENCES products(product_name) ON DELETE CASCADE,
	quantity_in_stock INT NOT NULL
);

CREATE OR REPLACE FUNCTION adjust_product_pricing()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products
    SET price = CEIL(price * 1.1) 
    WHERE product_name IN (
        SELECT product_name 
        FROM inventory 
        WHERE quantity_in_stock < 10
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER adjust_product_pricing_trigger
AFTER INSERT OR UPDATE ON inventory
FOR EACH ROW EXECUTE FUNCTION adjust_product_pricing();

