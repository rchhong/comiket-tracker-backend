# Comiket Display Backend

This repository contains the backend code for the Comiket Tracker Display, written in Go.
For now, the backend is hot-reloaded using [Air](https://github.com/air-verse/air).
This is only for development purposes.
Later on, this will be disabled.

## Set-up

1. Create a currency conversion API key here: [https://currency.getgeoapi.com/](https://currency.getgeoapi.com/)
2. Create `.env` with the following contents in the root of the project

```
POSTGRES_USER=<postgres_username>
POSTGRES_PASSWORD=<postgres_password>
POSTGRES_DB="comiket"
PGDATA="/var/lib/postgresql/17/docker"
CURRENCY_API_KEY=<currency_api_key>
```

3. Create `config.yaml` in `cmd/` with the following contents

```{yaml}
app:
  port: 3000

logging:
  logLevel: INFO
  file:
    logFilePath: "/app/logs/comiket_backend.log"

db:
  postgres:
    host: comiket-db
    port: 5432
    databaseName: comiket
    username: <postgres_username>
    password: <postgres_password>
```

```{bash}
docker-compose -f ./deployments/docker-compose.yml up --detatch
```

## Tear-down

```{bash}
docker-compose -f ./deployments/docker-compose.yml down
```

# Design

![](./assets/sql.png)
