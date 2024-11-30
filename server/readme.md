# how to setup the db stuff

1. connect to db
```
$ psql -U postgres
```

2. create the taskuser & database

some usefull commands
```sql
\l              -- list the database
\du             -- list the users
\q              -- leave
\c              -- connect to database
\c postgres     -- go back
\dt             -- show the tables
```

destroy the database:
```sql
drop database if exists area_database;
drop user if exists area_user;
```

```sql
create database area_database;
create user area_user with encrypted password 'password';
grant select, insert, update, delete on all tables in schema public to area_user;
-- grant all privileges on database area_database to area_user;
```

connect to the area_database
```sql
\c area_database;
```

create the table
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    encryptedPassword TEXT NOT NULL
);
```

if you want to change the password
```sql
ALTER USER area_user WITH PASSWORD 'newpassword';
```


## docker

pull image
```sh
docker pull postgres
```

run postregsql on docker
```sh
docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```

use psql console on container
```sh
docker exec -it postgres psql -U postgres
```

