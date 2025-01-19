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
    email TEXT default null,
    password TEXT default NULL,
    tokenid int default null
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

create table if not exists tokens (
    id serial primary key,
    service varchar(255) not null,
    token text not null,
    refresh text not null,
    -- tokenid text not null,
    -- userid int default null
    userid text not null,
    owner int default null
);

-- insert into users (tokenid) values (
--     1
-- );

-- insert into tokens (service, token, refresh, userid, owner) values (
--     'roblox',
--     '<token>',
--     '<refresh>',
--     '3462185362',
--     1
-- );


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
    ('Spotify', 'assets/services/spotify.png', '/service/spotify', '#05b348', '#038a2b'),
    ('Netflix', 'assets/services/netflix.png', '/service/netflix', '#e50914', '#b2070f'),
    ('Weather Underground', 'assets/services/weather-underground.webp', '/service/weather-underground', '#222222', '#000000'),
    ('Instagram', 'assets/services/instagram.webp', '/service/instagram', '#f56040', '#d84a2f'),
    ('Twitter', 'assets/services/x.webp', '/service/twitter', '#222222', '#000000'),
    ('Notification', 'assets/services/notification.webp', '/service/notification', '#222222', '#000000'),
    ('Android', 'assets/services/android.webp', '/service/android', '#3ddc84', '#2dbb6a'),
    ('Time', 'assets/services/time.webp', '/service/time', '#222222', '#000000'),
    ('Discord', 'assets/services/discord.webp', '/service/discord', '#7289da', '#5865f2'),
    ('Nasa', 'assets/services/nasa.webp', '/service/nasa', '#3d1f5e', '#2b1444');

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
    ('Create playlist of your favorite series in one click', '/applet/spotify/1', 132124, 1, 2),
    ('Get the weather forecast every day at 7:00 AM', '/applet/weather-underground/1', 88432, 3, 6),
    ('Update your Android wallpaper with NASA''s image of the day', 'applet/nasa/1', 348839, 8, 7),
    ('Tweet your Instagrams as native photos on Twitter', '/applet/instagram/1', 603723, 4, 5),
    ('Schedule sending of discord message', '/applet/discord/1', 324434, 9, 8);

