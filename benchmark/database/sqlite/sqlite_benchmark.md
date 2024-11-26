# SQlite Benchmark

<br>

# Sommaire
- [**Création d'un schéma SQL**](#1-création-dun-schéma-sql)
- [**Documentation sur l'installation**](#2-documentation-sur-linstallation)
- [**Exemples de requêtes SQL**](#3-exemples-de-requêtes-sql)
- [**Exemple de connexion**](#4-exemple-de-connexion)
- [**Points positifs et négatifs**](#5-analyse-des-points-positifs-et-négatifs)

<br>
<br>

# 1. Création d'un schéma SQL

## Schéma de base de données

cf benchmark.sql

## Résumé des contraintes appliquées dans cette base de données :

- **Clés primaires** : 
    - `id` dans chaque table (`categories`, `services`, `users`, `user_services`).
    - La combinaison des colonnes `user_id` et `service_id` dans la table `user_services` est définie comme clé primaire.

- **Clé étrangère** :
    - `category_id` dans la table `services`, reliée à `id` dans `categories`.
    - `user_id` dans la table `user_services`, reliée à `id` dans `users`.
    - `service_id` dans la table `user_services`, reliée à `id` dans `services`.

- **Contraintes uniques** :
    - Le champ `email` dans la table `users` est défini comme unique pour garantir qu'il n'y a pas de doublon.

<br>

# 2. Documentation sur l'installation

## Sommaire
- [**Installation sur Linux**](#installation-sur-linux)
- [**Installation sur MacOS**](#installation-sur-macos)
- [**Installation sur Windows**](#installation-sur-windows)

<br>

### Installation sur Linux

#### Mise à jour système :
```
sudo apt update
```

#### Installer MySQL :
```
sudo apt install sqlite3
```

### Charger la database :
```
sqlite3 benchmark.db < benchmark.sql
```

<br>

### Installation sur MacOS

#### Mise à jour système :
```
Auto update
```

#### Installer MySQL :
```
brew install sqlite
```

#### Charger la database :
```
sqlite3 benchmark.db < benchmark.sql
```

<br>

### Installation sur Windows

#### Installer un gestionnaire de paquets :
```
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

#### Installer MySQL :
```
choco install sqlite
```

#### Charger la database :
```
sqlite3 benchmark.db < benchmark.sql 
```

<br>

# 3. Exemples de requêtes SQL

### Insertion de données :
```sql
INSERT INTO categories (name, description) VALUES ('music', 'music'), ('video', 'video'), ('school', 'school');
```

### Lecture des données :
```sql
SELECT * FROM users;

SELECT users.name, users.surname, services.name AS service_name
FROM users
JOIN user_services ON users.id = user_services.user_id
JOIN services ON user_services.service_id = services.id;
```

### Mise à jour des données :
```sql
UPDATE users SET email = 'nouveau.email@example.com' WHERE id = 1;
```

### Suppression des données :
```sql
DELETE FROM user_services WHERE user_id = 1 AND service_id = 1;

DELETE FROM users WHERE id = 1;
```

<br>

# 4. Exemple de connexion

### Installer SQLite pour Python:
```
pip install sqlite3
```

### Exemple de script Python :

```py
import sqlite3

# Connexion à la base de données SQLite
conn = sqlite3.connect('benchmark.db')

cursor = conn.cursor()

# Exemple de requête : lire les utilisateurs
cursor.execute("SELECT * FROM users")
for row in cursor.fetchall():
    print(row)

# Fermer la connexion
conn.close()
```

# 5. Analyse des points positifs et négatifs


### Points positifs ✅ :
- **Légèreté** : Pas besoin d’un serveur, tout tient dans un fichier.

- **Facilité d’installation** : Prêt à l’emploi sur la plupart des systèmes.

- **Portabilité** : Une base SQLite est un simple fichier transférable.

- **Performance** : Idéal pour les petites et moyennes applications.

### Points négatifs ❌ :

- **Limitation de concurrence** : Moins performant pour les applications avec de nombreux utilisateurs simultanés.


- **Fonctionnalités limitées** : Pas de gestion avancée des transactions comme avec MySQL ou PostgreSQL.

- **Moins adapté pour les grands volumes de données** : Pas conçu pour des bases massives.
