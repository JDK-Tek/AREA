drop table if exists micro_github;

create table if not exists micro_github (
    id serial primary key,
    areauserid int not null,
    userid int not null,
    bridgeid int not null,
    triggers varchar(255) not null --,
);
