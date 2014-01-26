Perdocker
====

Evaluate code in different languages inside docker containers.

## Langs

Currently supported languages are:

- ruby (2.1.0)
- javascript (nodejs 0.10.24)
- golang (1.2)
- python (2.7.3)

## API

```bash
curl -POST -d "[1,2,3].each do |a| puts a*a; end;" 'http://localhost:8080/api/evaluate/ruby'
{"stdout":"1\n4\n9\n","stderr":"","exitCode":0}

curl -POST -d "var a = 6; a += 10; console.log(a)" 'http://localhost:8080/api/evaluate/nodejs'
{"stdout":"16\n","stderr":"","exitCode":0}

curl -POST -d "package main; import \"fmt\" ; func main() { fmt.Println(1+1) }" 'http://localhost:8080/api/evaluate/golang'
{"stdout":"2\n","stderr":"","exitCode":0}

curl -POST -d "print(\"Hello, World\")" 'http://localhost:8080/api/evaluate/python'
{"stdout":"Hello, World\n","stderr":"","exitCode":0}

curl -POST -d '{"language":"ruby", "code":"puts 1"}' 'http://localhost:8080/api/evaluate'
{"stdout":"1\n","stderr":"","exitCode":0}

```

## Install

```bash
make install
```

## Run

```bash
make run
```

> **NOTE:**
> Perdocker correctly works only with the latest dev Docker version. Caused by
> [this bug](https://github.com/dotcloud/docker/issues/1319). Bug
> will be fixed in 0.8.0 version. You can find compiled dev version
> of docker at `bin/` directory.
> You can install latest docker version and then just replace
> `/usr/bin/docker` with `./bin/docker` and restart service.

> **NOTE:**
> Perdocker expects that it will be able to run `docker` command without sudo.
> [Details](http://docs.docker.io/en/latest/use/basics/)

## Flags

```bash
./bin/perdocker -listen 127.0.0.1:80 -ruby-workers 5 -nodejs-workers 5 -golang-workers 5 -timeout 5
```

## Defaults

- listen :8080
- 1 ruby worker
- 1 nodejs worker
- 1 go worker
- 1 python worker
- timeout 30 seconds

## Coming soon

- timouts per eval.
- many language support (php, C, C++ and something else).
- improvement run process.
- start & attach to container instead of run it per request.
