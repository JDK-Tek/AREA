# You wanna create a microservice ?

First, create a folder in the `services` folder, with the name of your service.

Then make a web serer that listens on the port **80**.

This server can have the route `oauth` for oauth login, that can be
- `GET` to get the oauth link
- `POST` to register with oauth

When you register successfully with oauth, you shall send something like that:
```
200 OK
```
```json
{
    "token": "<the token>"
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

Generate a JWT that contains
- the userid `id` of the client.
- the expiration `exp` which is the token expiration.

Use the BACKEND_KEY environnement variable to encrypt the token.

Once everything is done, please add in the docker-compose the following rule:
```yml
  service-<service name>:
    env_file:
      - ./services/<service name>/.env
      - .env
    volumes:
      - ./services/<service name>/.env:/usr/mount.d/.env
      - ./.env:/usr/mount.d/.env1
    container_name: area-service-<service name>
    build:
      context: ./services/<service name>
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.<service name>.rule=PathPrefix(`/service/<service name>/`)"
      - "traefik.http.services.<service name>.loadbalancer.server.port=80"
      - "traefik.http.routers.<service name>.middlewares=stripprefix-service-<service name>"
      - "traefik.http.middlewares.stripprefix-service-<service name>.stripprefix.prefixes=/service/<service name>"
    environment:
      - DB_HOST=database
      - DB_PORT=${DATABASE_PORT}
      - DB_USER=${DATABASE_USER}
      - DB_PASSWORD=${DATABASE_PASSWORD}
      - DB_NAME=${DATABASE_NAME}
    networks:
      - web
```
Replace all `<service name>` by your actual service name.

Now it should work fine.
