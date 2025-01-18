drop table if exists micro_github;

create table if not exists micro_github (
    id serial primary key,
    userid int not null,
    bridgeid int not null,
    spices text not null,
    triggers varchar(255) not null --,
);
