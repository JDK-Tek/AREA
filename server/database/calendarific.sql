drop table if exists micro_calendarific;

create table if not exists micro_calendarific (
    id serial primary key,
    userid int not null,
    bridgeid int not null,
    spices text not null,
    events text not null,
    triggers varchar(255) not null --,
);
