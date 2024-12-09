-- drop database if exists area_database;
-- drop user if exists area_user;

-- create database area_database;
-- create user area_user with encrypted password 'password';
-- grant select, insert, update, delete on all tables in schema public to area_user;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS bridge;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

create table if not exists actions (
    id serial primary key,
    service varchar(255) not null,
    name varchar(255) not null,
    spices jsonb
);

create table if not exists reactions (
    id serial primary key,
    service varchar(255) not null,
    name varchar(255) not null,
    spices jsonb
);

create table if not exists bridge (
    id serial primary key,
    userid int,
    action int,
    reaction int --,
    -- foreign key (userid) references users(id) on delete cascade,
    -- foreign key (action) references actions(id) on delete cascade,
    -- foreign key (reaction) references reactions(id) on delete cascade
);