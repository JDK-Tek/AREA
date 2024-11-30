-- drop database if exists area_database;
-- drop user if exists area_user;

-- create database area_database;
-- create user area_user with encrypted password 'password';
-- grant select, insert, update, delete on all tables in schema public to area_user;

-- Create the users table if it does not already exist
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);
