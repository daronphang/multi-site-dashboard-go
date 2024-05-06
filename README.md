# Multi-site Dashboard Go

## Architecture

Web application is built using Clean Architecture by Uncle Bob.

https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

## First time setup

### Git

1. Clone git repository into your local directory

```sh
$ git clone https://tools.mf.platform/bitbucket/scm/dash/multi-site-dashboard-go.git
```

### Dependencies

1. Install dependencies

```sh
$ go mod download
```

### Wire

1. Setup wire for dependency injection

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

### sqlc

1. Install sqlc

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

5. Install statik for Go compilation into binary

```sh
$ go install github.com/rakyll/statik
```

6. Compile static files

```sh
$ statik -src path/to/root/directory/swaggerui
```

## Development

### Wire

https://github.com/google/wire

1. Write Providers and Wire functions

2. Generate Wire code

```sh
$ cd path/to/root/directory
$ wire internal
$ go generate internal # once wire_gen.go is created, can regenerate using this
```

### sqlc

https://docs.sqlc.dev/en/stable/

1. Write SQL queries

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
$ swagger serve ./docs/swagger.yaml --flavor swagger
```

## Deployment

### TimescaleDB

1. Deploy using PostgreSQL Kubernetes operators to simplify installation, configuration and lifecycle

https://github.com/zalando/postgres-operator/tree/master
