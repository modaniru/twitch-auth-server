CREATE TABLE users(
    id serial primary key,
    user_id varchar unique not null,
    client_id varchar unique not null,
    reg_date varchar not null
);