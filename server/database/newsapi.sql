drop table if exists micro_newsapi_articles;
drop table if exists micro_newsapi;

create table if not exists micro_newsapi_articles (
    id serial primary key,
    url text not null
);

create table if not exists micro_newsapi (
    id serial primary key,
    userid int not null,
    bridgeid int not null,
    spices text not null,
    triggers varchar(255) not null --,
);
