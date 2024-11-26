# Creation du benchmark mysql

<br>

# Sommaire
- [**Création d'un schéma SQL**](#1-création-dun-schéma-sql)
- [**Documentation sur l'installation**](#2-documentation-sur-linstallation)
- [**Exemples de requêtes SQL**](#3-exemples-de-requêtes-sql)
- [**Exemple de connexion**](#4-mini-example-de-connexion)
- [**Points positifs et négatifs**](#5-analyse-des-points-positifs-et-négatifs)

<br>
<br>

# 1. Création d’un schéma SQL

## Schéma de base de données

cf benchmark.sql

## Résumé des contraintes appliquées dans cette base de données :

- **Clés primaires** : id dans chaque table.

- **Clé étrangère** : utilisateur_id dans commandes, reliée à id dans utilisateurs.

- **Contraintes uniques** : champ email dans utilisateurs.

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
sudo apt update && sudo apt upgrade
```

#### Installer MySQL :
```
sudo apt install mysql-server
```

#### Démarrer le service :
```
sudo systemctl start mysql
```

#### Sécuriser l’installation :
```
sudo mysql_secure_installation
```

#### Charger la database :
```
mysql -u root -p benchmark < benchmark/database/mysql/benchmark.sql
```

#### Lancer le serveur :
```
mysql -u root -p benchmark
```

<br>

### Installation sur MacOS

#### Mise à jour système :
```
Auto update
```

#### Installer MySQL :
```
brew install mysql
brew install mysql@8.4
```

#### Démarrer le service :
```
brew services start mysql
```

#### Sécuriser l’installation :
```
mysql_secure_installation
```

#### Charger la database :
```
mysql -u root -p benchmark < benchmark/database/mysql/benchmark.sql
```

#### Lancer le serveur :
```
mysql -u root -p benchmark
```


<br>

### Installation sur Windows

#### Installer un gestionnaire de paquets :
```
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

#### Installer MySQL :
```
choco install mysql
```

#### Démarrer le service :
```
net start mysql
```

#### Sécuriser l’installation :
```
mysql_secure_installation
```

#### Charger la database :
```
mysql -u root -p benchmark < benchmark/database/mysql/benchmark.sql
```

#### Lancer le serveur :
```
mysql -u root -p benchmark
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

# 4. Mini example de connexion

### Installer le connecteur MySQL :

```
pip install mysql-connector-python
```

### Exemple de script :

```py
import mysql.connector

# Connexion à la base de données
conn = mysql.connector.connect(
    host="localhost",
    user="root",
    password="votre_mot_de_passe",
    database="GestionCommandes"
)

cursor = conn.cursor()

# Exécution d'une requête
cursor.execute("SELECT * FROM utilisateurs")
for row in cursor.fetchall():
    print(row)

conn.close()
```

# 5. Analyse des points positifs et négatifs

### Points positifs ✅ :
- **Facilité d’utilisation** : MySQL est bien documenté et bénéficie d’une large communauté.

- **Performance** : Bonne gestion des bases de données relationnelles de taille moyenne à grande.

- **Portabilité** : Fonctionne sur de nombreuses plateformes.

- **Gratuit** : Open source pour les besoins de base.

### Points négatifs ❌ :

- **Limites fonctionnelles** : Comparé à d'autres SGBD (ex. PostgreSQL), certaines fonctionnalités avancées sont moins développées.

- **Coûts supplémentaires** : Les fonctionnalités avancées peuvent nécessiter des licences (MySQL Enterprise).

- **Complexité pour les très grands volumes de données** : Performances moindres dans certains cas.

<br>
