create database product_management;

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

insert into products (name, price, quantity, created_at) VALUEs ('MacBook', 1500.50, 5, NOW());
insert into products (name, price, quantity, created_at) VALUEs ('iPhone', 999.99, 10, NOW());
insert into products (name, price, quantity, created_at) VALUEs ('Keyboard', 75.00, 20, NOW());

DELETE from products where id =1;

select *from products;

SELECT * FROM products WHERE price > 1000;

select *from products order by price asc;
                                             ^
select *from products order by price desc;

update products SET price = 1600.00 WHERE id =2;

update products set quantity = 15 WHERE id = 3;

delete from products where id =4;