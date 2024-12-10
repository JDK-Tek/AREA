## Launch my application

>[!NOTE]
> When i refer to `docker compose` command, you may want to use
> `docker-compose` on some versions of docker.

first, make sure you have `docker` & `docker compose` (or `docker-compose`)

then, make sure to have a `.env` file in this repository, that contains:
```py
# backend
BACKEND_PORT=...
BACKEND_KEY=...

# database
DB_PORT=...
DB_PASSWORD=...
DB_USER=...
DB_NAME=...
```

now you can run my app, from here, you can run:
```sh
docker compose up --build
```

this should run the app.

## Debugging

if you wanna debug the database, here is the command to connect to the docker:
```sh
docker exec -it area-database psql -U postgres -d area_database
```

you will be connected to the database, please see bellow to make postgres
commands


## How to get in touch with Postgres

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