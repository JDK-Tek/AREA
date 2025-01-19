drop table if exists micro_twitch_onstream;

create table if not exists micro_twitch_onstream (
    id serial primary key,
    streamer text not null,
    userid text not null,
    bridge int not null,
    areaid int not null,
    connected boolean not null default false
);
