# AREA Project

## Introduction 📚
This project is a multi-platform application developed to provide backend, mobile, and web services. It uses various technologies to offer a complete and integrated user experience.


## Technologies Used 🛠️
- **Backend**: GO
- **Database**: PostgreSQL
- **Mobile**: Flutter
- **Web**: ReactJS
- **Server**: Docker Compose, deployed on Azure using VMs

## Installation and Launch 🚀

### Prerequisites 📋
- .env file at the root of repository using this template:
```
COMPOSE_PATH_SEPARATOR=:
COMPOSE_FILE=docker-compose.yaml:server/docker-compose.yml

# web
WEB_PORT=8081
REACT_APP_BACKEND_URL=...

# general
REDIRECT=.../connected
FRONTEND=...
EXPIRATION=1800

# backend
BACKEND_PORT=8080
BACKEND_KEY=...
BACKEND_SERVICES_PATH=/usr/services

# database
DATABASE_PORT=42001
DATABASE_PASSWORD=...
DATABASE_USER=postgres
DATABASE_NAME=area_database

# traefik
TRAEFIK_PORT=42002
TRAEFIK_VERSION=2.10
```

- Docker and Docker Compose installed on your machine

### Using Docker 🐳
1. Clone the repository
2. Run `docker-compose up` in the root directory
3. The services will be available on the following ports (you can change them in the `.env` file):
    - Server: `localhost:8080`
    - Web: `localhost:8081`

## Deployment 🚀
The project is deployed on Azure using VMs.
<br>
The production web interface is available at [https://area.jepgo.root.sx/](https://area.jepgo.root.sx/) and the production server at [https://api.area.jepgo.root.sx/](https://api.area.jepgo.root.sx/).
<br>
The development web interface is available at [https://dev.area.jepgo.root.sx/](https://dev.area.jepgo.root.sx/) and the development server at [https://api.dev.area.jepgo.root.sx/](https://api.dev.area.jepgo.root.sx/).

## Authors ✨
- Elise PIPET
- Grégoire LANTIM
- Paul PARISOT
- Esteban MARQUES
- John DE KETTELBUTTER
