/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** techno
*/



const technoList = [
    {
        "name": "React Js",
        "type": "web",
        "img": '../assets/web/reactjs.png'
    },
    {
        "name": "HTML",
        "type": "web",
        "img": '../assets/web/html.png'
    },
    {
        "name": "VueJs",
        "type": "web",
        "img": '../assets/web/vuejs.png'
    },
    {
        "name": "Flutter",
        "type": "mobile",
        "img": '../assets/mobile/flutter.png'
    },
    {
        "name": "React Native",
        "type": "mobile",
        "img": '../assets/mobile/reactnative.png'
    },
    {
        "name": "Kotline",
        "type": "mobile",
        "img": '../assets/mobile/kotline.png'
    },
    {
        "name": "Django",
        "type": "backend",
        "img": '../assets/backend/django.png'
    },
    {
        "name": "Go",
        "type": "backend",
        "img": '../assets/backend/go.png'
    },
    {
        "name": "NestJs",
        "type": "backend",
        "img": '../assets/backend/nestjs.jpg'
    },
    {
        "name": "MariaDB",
        "type": "database",
        "img": '../assets/database/mariadb.jpeg'
    },
    {
        "name": "MySQL",
        "type": "database",
        "img": '../assets/database/mysql.png'
    },
    {
        "name": "PostgreSQL",
        "type": "database",
        "img": '../assets/database/postgres.png'
    },
    {
        "name": "SQLite",
        "type": "database",
        "img": '../assets/database/sqlite.png'
    }
];

const select = document.querySelector('select');
const divTechno = document.querySelector('#div-techno');
const onLoadLabel = document.querySelector('#load-time');

window.addEventListener('DOMContentLoaded', () => {
    select.dispatchEvent(new Event('change'));
});

select.addEventListener('change', (e) => {
    const startTime = performance.now();
    const selectedTechno = e.target.value;
    divTechno.innerHTML = '';
    technoList.forEach(techno => {
        if (selectedTechno === "all" || techno.type === selectedTechno) {
            const div = document.createElement('div');
            div.classList.add('card');
            div.innerHTML = `
                <img class="img-techno" src="${techno.img}" alt="${techno.name}">
            `;
            divTechno.appendChild(div);
        }
    });

    const endTime = performance.now();
    const loadTime = (endTime - startTime).toFixed(2);
    onLoadLabel.textContent = 'Temps de chargement: ' + loadTime + ' ms';
});
