# Creation of MySQL Benchmark

<br>

# Table of Contents
- [**Schema Creation**](#1-schema-creation)
- [**Installation Documentation**](#2-installation-documentation)
- [**SQL Query Examples**](#3-sql-query-examples)
- [**Connection Example**](#4-connection-example)
- [**Positive and Negative Points**](#5-analysis-of-positive-and-negative-points)

<br>
<br>

# 1. Schema Creation

## Database Schema

See `benchmark.sql`.

## Summary of constraints applied in this database:

- **Primary Keys**: `id` in each table.

- **Foreign Key**: `user_id` in `orders`, linked to `id` in `users`.

- **Unique Constraints**: `email` field in `users`.

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
sudo apt update && sudo apt upgrade
```

#### Install MySQL:
```bash
sudo apt install mysql-server
```

#### Start the service:
```bash
sudo systemctl start mysql
```

#### Secure the installation:
```bash
mysql_secure_installation --no-defaults
```
```plaintext
Would you like to setup VALIDATE PASSWORD component? y
Please enter 0 = LOW, 1 = MEDIUM, 2 = STRONG: 0
New password: **********
Re-enter password: **********

Do you wish to continue with the password provided? y
Remove anonymous users? y
Disallow root login remotely? n
Remove test database and access to it? y
Reload privilege tables now? y
```

#### Load the database:
```bash
mysql -u root -p < benchmark/database/mysql/benchmark.sql
Enter password: **********
```

#### Start the server:
```bash
mysql -u root -p benchmark
Enter password: **********
```


<br>

### Installation on MacOS

#### Update the system:
```bash
Auto update
```

#### Install MySQL:
```bash
brew install mysql
brew install mysql@8.4
```

#### Start the service:
```bash
brew services start mysql
```

#### Secure the installation:
```bash
mysql_secure_installation --no-defaults
```
```plaintext
Would you like to setup VALIDATE PASSWORD component? y
Please enter 0 = LOW, 1 = MEDIUM, 2 = STRONG: 0
New password: **********
Re-enter password: **********

Do you wish to continue with the password provided? y
Remove anonymous users? y
Disallow root login remotely? n
Remove test database and access to it? y
Reload privilege tables now? y
```

#### Load the database:
```bash
mysql -u root -p < benchmark/database/mysql/benchmark.sql
Enter password: **********
```

#### Start the server:
```bash
mysql -u root -p benchmark
Enter password: **********
```


<br>

### Installation on Windows

#### Install a package manager:
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

#### Install MySQL:
```powershell
choco install mysql
```

#### Start the service:
```powershell
net start mysql
```

#### Secure the installation:
```bash
mysql_secure_installation --no-defaults
```
```plaintext
Would you like to setup VALIDATE PASSWORD component? y
Please enter 0 = LOW, 1 = MEDIUM, 2 = STRONG: 0
New password: **********
Re-enter password: **********

Do you wish to continue with the password provided? y
Remove anonymous users? y
Disallow root login remotely? n
Remove test database and access to it? y
Reload privilege tables now? y
```

#### Load the database:
```bash
mysql -u root -p < benchmark/database/mysql/benchmark.sql
Enter password: **********
```

#### Start the server:
```bash
mysql -u root -p benchmark
Enter password: **********
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

### Install the MySQL connector:

```bash
pip install mysql-connector-python
```

### Example Script:

```python
import mysql.connector

# Connecting to the database
conn = mysql.connector.connect(
    host="localhost",
    port="3306",
    user="root",
    password="**********",
    database="benchmark"
)

cursor = conn.cursor()

# Execute a query
cursor.execute("SELECT * FROM users")
for row in cursor.fetchall():
    print(row)

conn.close()
```

# 5. Analysis of positive and negative points

### Positive Points ✅:
- **Ease of Use**: MySQL is well-documented and has a large community.

- **Performance**: Handles medium to large-scale relational databases efficiently.

- **Portability**: Works across multiple platforms.

- **Free**: Open-source for basic needs.

### Negative Points ❌:

- **Functional Limitations**: Compared to other DBMS (e.g., PostgreSQL), some advanced features are less developed.

- **Additional Costs**: Advanced features may require licenses (MySQL Enterprise).

- **Complexity for Very Large Data Volumes**: May face performance challenges in certain cases.
