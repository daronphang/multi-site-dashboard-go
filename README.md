# Multi-site Dashboard Go

## Architecture

Web application is built using Clean Architecture by Uncle Bob.

https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

## First time setup (local)

### Git

1. Clone git repository into your local directory

```sh
$ git clone https://tools.mf.platform/bitbucket/scm/dash/multi-site-dashboard-go.git
```

### Wire

1. Setup wire for dependency injection

https://github.com/google/wire

```sh
$ go install github.com/google/wire/cmd/wire@latest
$ go env # extract value of GOPATH
$ export GOPATH=<GOPATH> # add to .bashrc
$ export PATH=$PATH:$GOPATH/bin
$ source ~/.bashrc
```

2. Update VSCode to resolve import issues

```json
// command + shift + p to open User Settings
{
  "gopls": {
    "buildFlags": ["-tags=wireinject"]
  }
}
```

### TimescaleDB

1. Install as a container

https://docs.timescale.com/self-hosted/latest/install/installation-docker/

```sh
$ docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg16
```

2. Install psql (CLI tool) for PostgreSQL

```sh
$ brew install postgresql@<VERSION>
$ export PATH=$PATH:/opt/homebrew/opt/postgresql@<VERSION>/bin
$ psql --version
```

3. Install pgAdmin (GUI) for PostgreSQL

https://www.pgadmin.org/download/pgadmin-4-macos/

### Kafka

1. Install as a container

https://kafka.apache.org/quickstart

```sh
$ docker run --name kafka -p 9092:9092 -d apache/kafka:3.7.0
```

### sqlc

1. Install sqlc

https://github.com/sqlc-dev/sqlc

```sh
$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Swagger

https://ribice.medium.com/serve-swaggerui-within-your-golang-application-5486748a5ed4

1. Install Swagger

```sh
$ go install github.com/go-swagger/go-swagger/cmd/swagger
```

2. Clone swagger-ui source into any directory

```sh
$ cd path/to/any/directory
$ git clone https://github.com/swagger-api/swagger-ui.git
```

3. Copy dist folder from swagger-ui to swaggerui folder in source

```sh
cp -a path/to/swagger-ui/dist/ /path/to/root/directory/swaggerui
```

4. Change url in swagger-initializer.js

```
url: "./swagger.yaml"
```

### Statik

1. Install statik for Go compilation into binary (for swagger)

https://github.com/rakyll/statik

```sh
$ go install github.com/rakyll/statik
```

## Development

### Wire

https://github.com/google/wire

1. Write Providers and Wire functions

2. Generate Wire code

```sh
$ cd path/to/root/directory
$ wire ./internal
$ go generate ./internal # once wire_gen.go is created, can regenerate using this
```

### sqlc

https://docs.sqlc.dev/en/stable/

1. Provide SQL schemas and queries in `internal/database/migration` and `internal/repository/query` respectively

2. Generate code

```sh
$ cd path/to/root/directory
$ sqlc generate
```

### Web server

1. Run server

```sh
$ cd path/to/root/directory
$ go run cmd/rest/main.go
```

### Swagger

If serving using Echo, can navigate to `/api/v1/swagger/index.html`.

```sh
$ cd path/to/root/directory
$ swagger generate spec -o ./swaggerui/swagger.yaml --scan-models
$ statik -src path/to/root/directory/swaggerui # Rerun statik.
$ swagger serve ./docs/swagger.yaml --flavor swagger
```

## Testing

Before running tests, set environment variable GO_ENV to 'TESTING'.

```sh
$ export GO_ENV=TESTING
```

### Mockery

To generate mocks from interfaces, use Mockery.

https://github.com/vektra/mockery

```sh
$ brew install mockery
$ cd path/to/root/directory
$ mockery # reads from .mockery.yaml config file
```

### Running unit tests

```sh
$ cd path/to/root/directory
$ go test ./... -v
$ go test ./... -v -coverpkg=./...
```

### Running integration tests

1. Start containers with Docker Compose

```sh
$ cd path/to/root/directory
$ docker compose -f docker-compose-testing.yaml up -d
```

## Deployment

### AWS EC2

1. Clone repository into EC2

```sh
$ git clone <public_url>
```

2. Copy config file

```sh
$ cd path/to/root/directory
$ touch internal/config/config.production.yaml
```

3. Run Docker Compose

```sh
$ docker compose -f docker-compose-deployment.yaml up -d
```
