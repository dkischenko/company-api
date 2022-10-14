# company-api

XM Golang Exercise - v22.0.0

## How to run for local needs

```bash
docker run -it --rm \
--name some-postgres \
-e POSTGRES_PASSWORD=password \
-e POSTGRES_USER=postgresPOSTGRES_USER \
-p 5432:5432 \
-e PGDATA=/var/lib/postgresql/data/pgdata \
-v ~/local-go-postgres:/var/lib/postgresql/data \
postgres:14.0

export SIGNINKEY=

go run main.go
```

## Configuration

| Enviroment Variable | Description               | Default                                                                             |
| ------------- |:--------------------------|:------------------------------------------------------------------------------------|
| `HOST` | application host          | `127.0.0.1`                                                                         |
| `PORT` | application port          | `9090`                                                                              |
| `DATABASE_DSN` | Postgres database DSN     | `host=db user=postgres password=password dbname=postgres port=5432 sslmode=disable` |
| `ACCESS_TOKEN_TTL` | TTL of JWT token(seconds) | `120s`                                                                              |
| `SIGNINKEY` | Key to create signed JWT  | `10`                                                                                |