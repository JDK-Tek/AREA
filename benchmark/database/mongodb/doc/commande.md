# Shortcut - commands

## Load a database

```bash
mongosh /path/to/scheme/toload.js
```

## Displays

```js
show dbs // show all databases
use database // use an database
show tables // show tables of the database

db.dropDatabase() // drop the database
```

## Select rows

```js
db.users.findOne(); // list one user
db.users.find(); // list all users


db.users.find({ age: { $gt: 20 } }); // list all users age > 20
db.users.find({ age: { $gte: 20 } }); // list all users age >= 20
db.users.find({ age: { $eq: 20 } }); // list all users age == 20
db.users.find({ age: { $lte: 20 } }); // list all users age <= 20
db.users.find({ age: { $lt: 20 } }); // list all users age < 20

db.services.find({ category: db.categories.findOne({ name: 'video' })._id }); // list all services of category video
```


