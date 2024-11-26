# SQLite Benchmark

<br>

# Table of Contents
- [**Schema Creation**](#1-schema-creation)
- [**Installation Documentation**](#2-installation-documentation)
- [**SQL Query Examples**](#3-sql-query-examples)
- [**Connection Example**](#4-connection-example)
- [**Positive and Negative points**](#5-analysis-of-positive-and-negative-points)

<br>
<br>

# 1. Schema Creation

## Database Schema

See `benchmark.sql`.

## Summary of constraints applied in this database:

- **Primary Keys**: 
    - `id` in each table (`categories`, `services`, `users`, `user_services`).
    - The combination of `user_id` and `service_id` columns in the `user_services` table is defined as the primary key.

- **Foreign Keys**:
    - `category_id` in the `services` table, linked to `id` in `categories`.
    - `user_id` in the `user_services` table, linked to `id` in `users`.
    - `service_id` in the `user_services` table, linked to `id` in `services`.

- **Unique Constraints**:
    - The `email` field in the `users` table is defined as unique to ensure no duplicates.

<br>

# 2. Installation Documentation

## Table of Contents
- [**Installation on Linux**](#installation-on-linux)
- [**Installation on MacOS**](#installation-on-macos)
- [**Installation on Windows**](#installation-on-windows)

<br>

### Installation on Linux

#### Update the system:
```bash
sudo apt update
```

#### Install SQLite:
```bash
sudo apt install sqlite3
```

### Load the database:
```bash
sqlite3 benchmark.db < benchmark.sql
```

<br>

### Installation on MacOS

#### Update the system:
```bash
Auto update
```

#### Install SQLite:
```bash
brew install sqlite
```

#### Load the database:
```bash
sqlite3 benchmark.db < benchmark.sql
```

<br>

### Installation on Windows

#### Install a package manager:
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

#### Install SQLite:
```powershell
choco install sqlite
```

#### Load the database:
```bash
sqlite3 benchmark.db < benchmark.sql 
```

<br>

# 3. SQL Query Examples

### Data Insertion:
```sql
INSERT INTO categories (name, description) VALUES ('music', 'music'), ('video', 'video'), ('school', 'school');
```

### Data Retrieval:
```sql
SELECT * FROM users;

SELECT users.name, users.surname, services.name AS service_name
FROM users
JOIN user_services ON users.id = user_services.user_id
JOIN services ON user_services.service_id = services.id;
```

### Data Update:
```sql
UPDATE users SET email = 'new.email@example.com' WHERE id = 1;
```

### Data Deletion:
```sql
DELETE FROM user_services WHERE user_id = 1 AND service_id = 1;

DELETE FROM users WHERE id = 1;
```

<br>

# 4. Connection Example

### Install SQLite for Python:
```bash
pip install sqlite3
```

### Python Script Example:

```python
import sqlite3

# Connect to the SQLite database
conn = sqlite3.connect('benchmark.db')

cursor = conn.cursor()

# Example query: read users
cursor.execute("SELECT * FROM users")
for row in cursor.fetchall():
    print(row)

# Close the connection
conn.close()
```

# 5. Analysis of positive and negative points

### Positive Points ✅:
- **Lightweight**: No server required, everything is stored in a single file.

- **Easy to install**: Ready to use on most systems.

- **Portability**: An SQLite database is a simple, transferable file.

- **Performance**: Ideal for small to medium applications.

### Negative Points ❌:

- **Concurrency limitations**: Less efficient for applications with many simultaneous users.

- **Limited features**: Lacks advanced transaction management like MySQL or PostgreSQL.

- **Less suited for large datasets**: Not designed for massive databases.
