/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** banchmark
*/

// db.Benchmark.dropDatabase();

const categories = [
    { name: 'music', description: 'music' },
    { name: 'video', description: 'video' },
    { name: 'school', description: 'school' }
];

const categoriesIds = [];
categories.forEach((category) => {
    const result = db.categories.insertOne(category);
    categoriesIds.push(result.insertedId);
});

const services = [
    { name: 'spotify', description: 'music streaming', url: 'https://www.spotify.com', category: categoriesIds[0] },
    { name: 'netflix', description: 'video streaming', url: 'https://www.netflix.com', category: categoriesIds[1] },
    { name: 'epitech', description: 'school', url: 'https://www.epitech.eu', category: categoriesIds[2] },
    { name: 'youtube', description: 'video streaming', url: 'https://www.youtube.com', category: categoriesIds[1] },
    { name: 'twitch', description: 'video streaming', url: 'https://www.twitch.com', category: categoriesIds[1] }
];

const servicesIds = [];
services.forEach((service) => {
    const result = db.services.insertOne(service);
    servicesIds.push(result.insertedId);
});

db.users.insertMany([
    {name: 'esteban', surname: 'marques', age: 11, email: 'esteban.marques@epitech.eu', password: '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', services: [servicesIds[0], servicesIds[1]]},
    {name: 'john', surname: 'de kettelbutter', age: 20, email: 'john.de-kettelbutter@epitech.eu', password: '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', services: [servicesIds[2], servicesIds[3]]},
    {name: 'paul', surname: 'parisot', age: 3, email: 'paul.parisot@epitech.eu', password: '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', services: [servicesIds[4]]},
    {name: 'elise', surname: 'pipet', age: 21, email: 'elise.pipet@epitech.eu', password: '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', services: [servicesIds[0], servicesIds[1], servicesIds[2]]},
    {name: 'gregoire', surname: 'lan tim', age: 22, email: 'gregpire.lan-tim@epitech.eu', password: '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', services: [servicesIds[3], servicesIds[4]]},
]);

