# Multi-site Dashboard Go

## First time setup

1. Clone git repository

```sh
$ git clone
```

2. Setup wire for dependency injection

```sh
$ go install github.com/google/wire/cmd/wire@latest
$ go env # extract value of GOPATH
$ export GOPATH=<GOPATH> # add to .bashrc
$ export PATH=$PATH:$GOPATH/bin
$ source ~/.bashrc
```

3. Update VSCode to resolve import issues when using wire

```json
// command + shift + p to open User Settings
{
  "gopls": {
    "buildFlags": ["-tags=wireinject"]
  }
}
```

4. Generate wire injectors

```sh
$ cd path/to/root/directory
$ wire
$ go generate # once wire_gen.go is created, can regenerate using this
```

5. Running application

```sh
$ go run main.go wire_gen.go
```
