drop table if exists micro_openweather;

create table if not exists micro_openweather (
    id serial primary key,
    userid int not null,
    bridgeid int not null,
    spices text not null,
    triggers varchar(255) not null,
    last_weather varchar(255) not null
);
