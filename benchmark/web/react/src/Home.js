import React, { useState } from 'react';
import BasicMenu from './BasicMenu'

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
  
    const renderImages = () => {
      switch (selectedCategory) {
        case 'Backend':
          return (
            <div>
              <img src={djangoLogo} alt="Django Logo" />
              <img src={goLogo} alt="Go Logo" />
              <img src={nestjsLogo} alt="NestJS Logo" />
            </div>
          );
        case 'Database':
          return (
            <div>
              <img src={mariadbLogo} alt="MariaDB Logo" />
              <img src={mysqlLogo} alt="MySQL Logo" />
              <img src={postgresLogo} alt="Postgres Logo" />
              <img src={sqliteLogo} alt="SQLite Logo" />
            </div>
          );
        case 'Web':
          return (
            <div>
              <img src={htmlLogo} alt="HTML Logo" />
              <img src={reactjsLogo} alt="ReactJS Logo" />
              <img src={vuejsLogo} alt="VueJS Logo" />
            </div>
          );
        case 'Mobile':
          return (
            <div>
              <img src={flutterLogo} alt="Flutter Logo" />
              <img src={kotlinLogo} alt="Kotlin Logo" />
              <img src={reactnativeLogo} alt="React Native Logo" />
            </div>
          );
        default:
          return null;
      }
    };
  
    return (
      <div>
        <BasicMenu onClick={handleMenuClick} />
        {renderImages()}
      </div>
    );
  }

export default Home;