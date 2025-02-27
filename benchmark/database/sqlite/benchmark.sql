CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- Clé primaire pour chaque catégorie
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS services (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- Clé primaire pour chaque service
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    url VARCHAR(255) NOT NULL,
    category_id INT,  -- Clé étrangère reliant la catégorie
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- Clé primaire pour chaque utilisateur
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_services (
    user_id INT,  -- Clé étrangère vers les utilisateurs
    service_id INT,  -- Clé étrangère vers les services
    PRIMARY KEY (user_id, service_id),  -- La combinaison des deux clés est unique
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO categories (name, description) VALUES 
('music', 'music'),
('video', 'video'),
('school', 'school');

INSERT OR IGNORE INTO services (name, description, url, category_id) VALUES 
('spotify', 'music streaming', 'https://www.spotify.com', 1),
('netflix', 'video streaming', 'https://www.netflix.com', 2),
('epitech', 'school', 'https://www.epitech.eu', 3),
('youtube', 'video streaming', 'https://www.youtube.com', 2),
('twitch', 'video streaming', 'https://www.twitch.com', 2);

INSERT OR IGNORE INTO users (name, surname, age, email, password) VALUES 
('esteban', 'marques', 11, 'esteban.marques@epitech.eu', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
('john', 'de kettelbutter', 20, 'john.de-kettelbutter@epitech.eu', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
('paul', 'parisot', 3, 'paul.parisot@epitech.eu', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
('elise', 'pipet', 21, 'elise.pipet@epitech.eu', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
('gregoire', 'lan tim', 22, 'gregpire.lan-tim@epitech.eu', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');

INSERT OR IGNORE INTO user_services (user_id, service_id) VALUES
(1, 1), (1, 2),
(2, 3), (2, 4),
(3, 5),
(4, 1), (4, 2), (4, 3),
(5, 4), (5, 5);
