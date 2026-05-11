CREATE TABLE products (
    id UUID PRIMARY KEY,
    count INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE reservations (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id),
    count INTEGER NOT NULL
);
