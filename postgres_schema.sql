-- customer table
create table customers (
    id serial primary key,
    customer_id varchar(50) unique not null,
    name varchar(100),
    email varchar(50),
    address varchar(500),
    created_by varchar(50),
    created_date timestamptz,
    updated_by varchar(50),
    updated_date timestamptz 
);

-- orders table
create table orders (
    id serial primary key,
    order_id varchar(100) unique not null,
    customer_id varchar(50),
    region varchar(100),
    date_of_sale date,
    payment_method varchar(250),
    shipping_cost numeric(10,2),
    created_by varchar(50),
    created_date timestamptz,
    updated_by varchar(50),
    updated_date timestamptz,
    foreign key (customer_id) references customers(customer_id)
);


-- product table
create table products (
    id serial primary key,
    product_id varchar(50) unique not null,
    product_name varchar(100),
    category varchar(100),
    created_by varchar(50),
    created_date timestamptz,
    updated_by varchar(50),
    updated_date timestamptz 
);

-- order items table
create table order_items (
    id serial primary key,
    order_id varchar(100),
    product_id varchar(50),
    quantity_sold int,
    unit_price numeric(10,2),
    discount numeric(5,2),
    created_by varchar(50),
    created_date timestamptz,
    updated_by varchar(50),
    updated_date timestamptz,
    foreign key (order_id) references orders(order_id) on delete cascade,
    foreign key (product_id) references products(product_id) on delete cascade
);
