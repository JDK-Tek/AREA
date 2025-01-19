# AREA Routes

## `POST` register
> /api/register

### Body
```json
{
    "email": "<email>",
    "password": "<password>"
}
```

## `POST` login
> /api/login

### Body
```json
{
    "email": "<email>",
    "password": "<password>"
}
```

## `POST` new area
> /api/area

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

## `GET` oauth getter
> /api/oauth/[service]

Get the oauth for a service.

## `POST` oauth setter
> /api/oauth/[service]

Set the oauth result (token, code to get token...)

### Body
```json
{
    "code": "<code>"
}
```

It returns a session token if it succeed, or anything if its an error.

## `GET` applets
> /api/area

Get the user applets, or the example applets.

If you put a token, it should return your applets.

However, if no token is found, it will 

### Request Headers (optional)
| Key | Value |
|-----|-------|
| Authorization | Bearer `<token>` |

### Query Params (optional)
| Key | Type |
|-----|-------|
| limit | `<number>` |

## `GET` doctor

> /api/doctor

Gives you usefull inform

### Request Headers (optional)
| Key | Value |
|-----|-------|
| Authorization | Bearer `<token>` |

---

## `PUT` orchestrator
> /api/orchestrator

>[!NOTE]
> This shouldnt be call manually

### Body
```json
{
    "bridge": <number>,
    "userid": <number>,
    "userid": {
        <string>: <string>
    },
}
```

