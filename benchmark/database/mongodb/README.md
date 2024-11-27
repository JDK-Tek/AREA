### 📊 **MongoDB Benchmark**

MongoDB is a NoSQL database, ideal for applications requiring horizontal scalability, fast writes, and schema flexibility.

Here’s a quick benchmark of MongoDB with its advantages, disadvantages, and how to use it with a backend.


## 🚀 **Installing MongoDB**

### On Linux (Ubuntu/Debian)

1. **Add the repository**:
   ```bash
   sudo apt update
   sudo apt install -y wget
   wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo tee /etc/apt/trusted.gpg.d/mongodb.asc
   echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -sc)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
   sudo apt update
   ```

2. **Install MongoDB**:
   ```bash
   sudo apt install -y mongodb-org
   ```

3. **Start MongoDB**:
   ```bash
   sudo systemctl start mongod
   sudo systemctl enable mongod
   ```

### On macOS (via Homebrew)

```bash
brew tap mongodb/brew
brew install mongodb-community@6.0
brew services start mongodb/brew/mongodb-community
```


## 🖥️ **Deploying MongoDB**

### **Local Deployment**

- **Simple and quick**: A MongoDB instance on a local server, ideal for small applications or testing purposes.

### **Production Deployment**

- **Replica Sets**: High availability with a primary node and secondary replicas.
- **Sharding**: Data is distributed across multiple servers to handle large-scale datasets.


## ✅❌ **Advantages and Disadvantages of MongoDB**

| **Advantages**                                                                | **Disadvantages**                                                                |
|-------------------------------------------------------------------------------|----------------------------------------------------------------------------------|
| ⚡ **Horizontal Scalability**: Easy to scale across multiple servers.          | 🔄 **Limited Transactions**: No complex ACID transactions until v4.0.            |
| 🗂️ **Flexible Data**: Stores JSON documents with dynamic schemas.              | 🧠 **Memory Usage**: Can be high for large databases.                            |
| 💪 **High Performance** for both writes and reads.                            | 🚫 **No Complex Joins**: Less efficient than relational databases.               |
| 🔒 **High Availability** with Replica Sets for better fault tolerance.         | ⚙️ **No Rigid Schema**: Can lead to inconsistencies in data.                     |


## ⚙️ **Using MongoDB with a Backend**

### **Node.js and Mongoose**

#### Installing Mongoose

```bash
npm install mongoose
```

#### Connecting to MongoDB

```javascript
const mongoose = require('mongoose');

mongoose.connect('mongodb://localhost:27017/my_database', { useNewUrlParser: true, useUnifiedTopology: true })
  .then(() => console.log("Connected to MongoDB!"))
  .catch(err => console.error("Connection error:", err));
```

#### Creating a Model and Saving a User

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
  .then(() => console.log('User saved!'))
  .catch(err => console.error('Error:', err));
```

#### Querying Users with Age Greater Than 20

```javascript
User.find({ age: { $gt: 20 } })
  .then(users => console.log('Users older than 20:', users))
  .catch(err => console.error('Error:', err));
```


## 💡 **Conclusion**

MongoDB is an excellent choice for applications requiring flexible data management, horizontal scalability, and high performance. However, for applications needing transactions or complex joins, it may have limitations.

Thus, in the context of our AREA project, MongoDB might be challenging to implement.