CREATE TABLE IF NOT EXISTS menu (
    id varchar(128) primary key,
    name varchar(1000),
    unitPrice int,
    currency varchar(10)
);