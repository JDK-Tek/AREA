drop table if exists micro_time;

create table if not exists micro_time (
    id serial primary key,
    userid int not null,
    bridgeid int not null,
    triggers timestamp not null,
    original timestamp not null --,
    -- foreign key (bridgeid) references bridge(id) on delete cascade
);
