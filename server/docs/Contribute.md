# You wanna create a microservice ?

First, create a folder in the `services` folder, with the name of your service.

Then make a web serer that listens on the port **80**.

## OAUTH

This server can have the route `oauth` for oauth login, that can be
- `GET` to get the oauth link
- `POST` to register with oauth

When you register successfully with oauth, you shall send something like that:
```
200 OK
```
```json
{
    "token": "<the token>" // (1)
}
```

Otherwise:
```
4xx ...
```
```json
{
    "error": "<why>"
}
```

`(1)` Generate a JWT that contains
- the userid `id` of the client.
- the expiration `exp` which is the token expiration.

Use the BACKEND_KEY environnement variable to encrypt the token.

## Deployment

Once everything is done, please add in the docker-compose the following rule:
```yml
  service-<your-service-name>:
    env_file:
      - .env
    container_name: area-service-<your-service-name>
    build:
      context: ./services/<your-service-name>
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.<your-service-name>.rule=PathPrefix(`/service/<your-service-name>/`)"
      - "traefik.http.services.<your-service-name>.loadbalancer.server.port=80"
      - "traefik.http.routers.<your-service-name>.middlewares=stripprefix-service-<your-service-name>"
      - "traefik.http.middlewares.stripprefix-service-<your-service-name>.stripprefix.prefixes=/service/<your-service-name>"
    environment:
      - BACKEND_KEY=${BACKEND_KEY}
      - REDIRECT=${REDIRECT}
      - DB_HOST=database
      - DB_PORT=${DATABASE_PORT}
      - DB_USER=${DATABASE_USER}
      - DB_PASSWORD=${DATABASE_PASSWORD}
      - DB_NAME=${DATABASE_NAME}
    networks:
      - web
```
Replace all `<your-service-name>` by your actual service name.

Now it should work fine.
