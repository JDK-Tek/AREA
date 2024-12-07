# AREA Routes

## `POST` register
> http://localhost:42000/api/register

### Body
```json
{
    "email": "<email>",
    "password": "<password>"
}
```

## `POST` login
> http://localhost:42000/api/login

### Body
```json
{
    "email": "<email>",
    "password": "<password>"
}
```

## `POST` new area
> http://localhost:42000/api/area

### Request Headers
| Key | Value |
|-----|-------|
| Authorization | Bearer `<token>` |

### Body
```json
{
    "action": {
        "service": "<service>",
        "name": "<reaction name>",
        "spices": {
            ...
        }
    },
    "reaction": {
        "service": "<service>",
        "name": "<reaction name>",
        "spices": {
            ...
        }
    }
}
```

### Example
```json
{
    "action": {
        "service": "time",
        "name": "in",
        "spices": {
            "howmuch": 2,
            "unit": "minutes"
        }
    },
    "reaction": {
        "service": "discord",
        "name": "send",
        "spices": {
            "channel": 0000000000000,
            "message": "Hello world !"
        }
    }
}
```

## `PUT` orchestrator
> http://localhost:42000/api/orchestrator

>[!NOTE]
> This shouldnt be call manually

### Body
```json
{
    "bridge": <number>
}
```

---

>[!NOTE]
> All the following bellow are the microservices, and should not be call manually

## `POST` discord send
> http://localhost:42002/service/discord/send

```json
{
    "spices": {
        ...
    }
}
```

## `POST` time in
> http://localhost:42002/service/time/in

```json
{
    "spices": {
        ...
    },
    "bridge": <number>
}
```

## `GET` discord oauth
> http://localhost:42002/service/discord/oauth

## Query Params
| Key | Value |
|-----|-------|
| redirect | <link>|
