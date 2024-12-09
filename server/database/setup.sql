-- drop database if exists area_database;
-- drop user if exists area_user;

-- create database area_database;
-- create user area_user with encrypted password 'password';
-- grant select, insert, update, delete on all tables in schema public to area_user;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS bridge;
drop table if exists token;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    spices JSONB
);

CREATE TABLE IF NOT EXISTS reactions (
    id SERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    spices JSONB
);

CREATE TABLE IF NOT EXISTS bridge (
    id SERIAL PRIMARY KEY,
    userid INT,
    action INT,
    reaction INT --,
    -- foreign key (userid) references users(id) on delete cascade,
    -- foreign key (action) references actions(id) on delete cascade,
    -- foreign key (reaction) references reactions(id) on delete cascade
);

create table if not exists token (
    id serial primary key,
    service varchar(255) not null,
    token text not null,
    userid int not null
);


-- EXEMPLES FOR MVP

DROP TABLE IF EXISTS applets;
DROP TABLE IF EXISTS services;

CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logo VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    colorN VARCHAR(7) NOT NULL,
    colorH VARCHAR(7) NOT NULL
);

INSERT INTO services (name, logo, link, colorN, colorH)
VALUES
    ('Spotify', './assets/services/spotify.png', 'https://www.spotify.com', '#05b348', '#038a2b'),
    ('Netflix', './assets/services/netflix.png', 'https://www.netflix.com', '#e50914', '#b2070f'),
    ('Weather Underground', './assets/services/weather-underground.png', 'https://www.wunderground.com/', '#222222', '#000000'),
    ('Instagram', './assets/services/instagram.png', 'https://www.instagram.com', '#f56040', '#d84a2f'),
    ('Twitter', './assets/services/x.png', 'https://www.x.com', '#222222', '#000000'),
    ('Notification', './assets/services/notification.png', '/notification', '#222222', '#000000'),
    ('Android', './assets/services/android.png', 'https://www.android.com', '#3ddc84', '#2dbb6a'),
    ('Nasa', './assets/services/nasa.png', 'https://www.nasa.com', '#3d1f5e', '#2b1444');

CREATE TABLE IF NOT EXISTS applets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    service1 INT NOT NULL,
    service2 INT NOT NULL,
    users INT,
    link VARCHAR(255) NOT NULL,
    FOREIGN KEY (service1) REFERENCES services(id) ON DELETE CASCADE,
    FOREIGN KEY (service2) REFERENCES services(id) ON DELETE CASCADE
);

INSERT INTO applets (name, link, users, service1, service2)
VALUES
    ('Create playlist of your favorite series in one click', 'https://spotify.com', 132124, 1, 2),
    ('Get the weather forecast every day at 7:00 AM', 'https://www.wunderground.com/', 88432, 3, 6),
    ('Update your Android wallpaper with NASA''s image of the day', 'https://www.nasa.gov/', 348839, 8, 7),
    ('Tweet your Instagrams as native photos on Twitter', 'https://instagram.com', 603723, 4, 5);

