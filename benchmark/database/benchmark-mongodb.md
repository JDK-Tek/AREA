### 📊 **Benchmark de MongoDB**

MongoDB est une base de données NoSQL, idéale pour les applications nécessitant une scalabilité horizontale, des écritures rapides, et une flexibilité de schéma.

Voici un benchmark rapide de MongoDB avec ses avantages, inconvénients, et comment l'utiliser avec un backend.

## 🚀 **Installation de MongoDB**

### Sur Linux (Ubuntu/Debian)

1. **Ajouter le dépôt** :
   ```bash
   sudo apt update
   sudo apt install -y wget
   wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo tee /etc/apt/trusted.gpg.d/mongodb.asc
   echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -sc)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
   sudo apt update
   ```

2. **Installer MongoDB** :
   ```bash
   sudo apt install -y mongodb-org
   ```

3. **Lancer MongoDB** :
   ```bash
   sudo systemctl start mongod
   sudo systemctl enable mongod
   ```

### Sur macOS (via Homebrew)

```bash
brew tap mongodb/brew
brew install mongodb-community@6.0
brew services start mongodb/brew/mongodb-community
```


## 🖥️ **Déploiement de MongoDB**

### **Déploiement Local**

- **Simple et rapide** : Une instance MongoDB sur un serveur local, idéale pour les petites applications ou les tests.

### **Déploiement en Production**

- **Replica Sets** : Haute disponibilité avec un nœud primaire et des répliques secondaires.
- **Sharding** : Répartition des données entre plusieurs serveurs pour gérer de grandes quantités de données.


## ✅❌ **Avantages et Inconvénients de MongoDB**

| **Avantages**                                                                 | **Inconvénients**                                                                |
|-------------------------------------------------------------------------------|----------------------------------------------------------------------------------|
| ⚡ **Scalabilité horizontale** : Facile à étendre sur plusieurs serveurs.      | 🔄 **Transactions limitées** : Pas de transactions ACID complexes jusqu'à la v4.0. |
| 🗂️ **Flexibilité des données** : Stocke des documents JSON avec des schémas dynamiques. | 🧠 **Consommation de mémoire** : Peut être élevé pour les grandes bases de données. |
| 💪 **Performances élevées** pour les écritures et les lectures.               | 🚫 **Pas de jointures complexes** : Moins performantes que dans les bases relationnelles. |
| 🔒 **Haute disponibilité** avec Replica Sets pour une meilleure tolérance aux pannes. | ⚙️ **Pas de schéma rigide** : Peut entraîner des incohérences dans les données.   |

## ⚙️ **Utilisation de MongoDB avec un Backend**

### **Node.js et Mongoose**

#### Installation de Mongoose

```bash
npm install mongoose
```

#### Connexion à MongoDB

```javascript
const mongoose = require('mongoose');

mongoose.connect('mongodb://localhost:27017/mon_database', { useNewUrlParser: true, useUnifiedTopology: true })
  .then(() => console.log("Connecté à MongoDB !"))
  .catch(err => console.error("Erreur de connexion:", err));
```

#### Créer un Modèle et Sauvegarder un Utilisateur

```javascript
const userSchema = new mongoose.Schema({
  name: String,
  age: Number,
  email: String
});

const User = mongoose.model('User', userSchema);

const newUser = new User({
  name: 'John Doe',
  age: 25,
  email: 'johndoe@example.com'
});

newUser.save()
  .then(() => console.log('Utilisateur sauvegardé !'))
  .catch(err => console.error('Erreur :', err));
```

#### Requête pour Trouver les Utilisateurs avec un Age Supérieur à 20

```javascript
User.find({ age: { $gt: 20 } })
  .then(users => console.log('Utilisateurs plus vieux que 20:', users))
  .catch(err => console.error('Erreur:', err));
```

## 💡 **Conclusion**

MongoDB est un excellent choix pour des applications nécessitant une gestion flexible des données, une scalabilité horizontale et des performances élevées. Cependant, pour des applications nécessitant des transactions ou des jointures, il peut avoir des limites.

Ainsi dans le cadre de notre projet AREA, mongodb risque d'être compliqué à mettre en place.