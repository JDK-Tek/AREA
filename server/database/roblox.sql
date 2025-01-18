drop database if exists micro_roblox;

create table if not exists micro_roblox (
    id serial primary key,
    robloxid text not null,
    gameid int not null,
    command text default null,
    constraint unique_gameid unique (gameid)
);
