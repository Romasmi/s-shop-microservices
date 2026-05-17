CREATE TABLE couriers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE slots (
    order_id UUID PRIMARY KEY,
    courier_id UUID NOT NULL REFERENCES couriers(id),
    from_date BIGINT NOT NULL,
    to_date BIGINT NOT NULL
);
