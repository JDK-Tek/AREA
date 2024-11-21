import React, { useState } from 'react';
import BasicMenu from './BasicMenu'
import './App.css'

// Backend logos
import djangoLogo from './assets/backend/django.png'
import goLogo from './assets/backend/go.png'
import nestjsLogo from './assets/backend/nestjs.jpg'

// Database logos
import mariadbLogo from './assets/database/mariadb.jpeg'
import mysqlLogo from './assets/database/mysql.png'
import postgresLogo from './assets/database/postgres.png'
import sqliteLogo from './assets/database/sqlite.png'

// Web logos
import htmlLogo from './assets/web/html.png'
import reactjsLogo from './assets/web/reactjs.png'
import vuejsLogo from './assets/web/vuejs.png'

// Mobile logos
import flutterLogo from './assets/mobile/flutter.png'
import kotlinLogo from './assets/mobile/kotline.png'
import reactnativeLogo from './assets/mobile/reactnative.png'

function Home() {
    const [selectedCategory, setSelectedCategory] = useState('Web');

    const handleMenuClick = (category) => {
        setSelectedCategory(category);
    };

    const images = {
        Backend: [
            { src: djangoLogo, alt: 'Django Logo' },
            { src: goLogo, alt: 'Go Logo' },
            { src: nestjsLogo, alt: 'NestJS Logo' }
        ],
        Database: [
            { src: mariadbLogo, alt: 'MariaDB Logo' },
            { src: mysqlLogo, alt: 'MySQL Logo' },
            { src: postgresLogo, alt: 'Postgres Logo' },
            { src: sqliteLogo, alt: 'SQLite Logo' }
        ],
        Web: [
            { src: htmlLogo, alt: 'HTML Logo' },
            { src: reactjsLogo, alt: 'ReactJS Logo' },
            { src: vuejsLogo, alt: 'VueJS Logo' }
        ],
        Mobile: [
            { src: flutterLogo, alt: 'Flutter Logo' },
            { src: kotlinLogo, alt: 'Kotlin Logo' },
            { src: reactnativeLogo, alt: 'React Native Logo' }
        ]
    };

    const renderImages = () => {
        const categoryImages = images[selectedCategory] || [];
        return (
            <div class="div-techno">
                {categoryImages.map((image, index) => (
                    <img class="img-techno" key={index} src={image.src} alt={image.alt} />
                ))}
            </div>
        );
    };

    return (
        <div>
            <BasicMenu onClick={handleMenuClick} />
            {renderImages()}
        </div>
    );
}

export default Home;