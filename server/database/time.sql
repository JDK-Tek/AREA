drop table if exists micro_time;

create table if not exists micro_time (
    id serial primary key,
    bridgeid int not null,
    triggers timestamp not null --,
    -- foreign key (bridgeid) references bridge(id) on delete cascade
);
