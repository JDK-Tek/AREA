drop table if exists micro_reddit;

create table if not exists micro_reddit (
    id serial primary key,
    bridgeid int not null,
    triggers timestamp not null
);
