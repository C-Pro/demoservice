create table users(id bigserial primary key, name varchar(128) unique not null, password_hash varchar(1024));
