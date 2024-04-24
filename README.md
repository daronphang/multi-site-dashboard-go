# Multi-site Dashboard Go

## First time setup

### Git

1. Clone git repository into your local directory

```sh
$ git clone https://tools.mf.platform/bitbucket/scm/dash/multi-site-dashboard-go.git
```

### Wire (dependency injection)

1. Setup wire

```sh
$ go install github.com/google/wire/cmd/wire@latest
$ go env # extract value of GOPATH
$ export GOPATH=<GOPATH> # add to .bashrc
$ export PATH=$PATH:$GOPATH/bin
$ source ~/.bashrc
```

3. Update VSCode to resolve import issues

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

2. Install psql (command line) for PostgreSQL

```sh
$ brew install postgresql@<VERSION>
$ export PATH=$PATH:/opt/homebrew/opt/postgresql@<VERSION>/bin
$ psql --version
```

3. Install pgAdmin (GUI) for PostgreSQL

https://www.pgadmin.org/download/pgadmin-4-macos/

## Development

### Wire

1. Generate code

```sh
$ cd path/to/root/directory
$ wire internal
$ go generate internal # once wire_gen.go is created, can regenerate using this
```

### SQL queries (sqlc)

1. Write queries in SQL (follow sqlc for documentation)

2. Generate code

### Web server

1. Run web server

```sh
$ go run cmd/main.go
```

## Deployment

### TimescaleDB

1. Deploy using PostgreSQL Kubernetes operators to simplify installation, configuration and lifecycle

https://github.com/zalando/postgres-operator/tree/master
