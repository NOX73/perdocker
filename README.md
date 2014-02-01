Perdocker
====

Evaluate code in different languages inside docker containers.

### Status

This project steel under active development. So feel free to open issue
if your have a problem.

## Langs

Currently supported languages are:

- ruby (2.1.0)
- javascript (nodejs 0.10.24)
- golang (1.2)
- python (2.7.3)
- C (gcc 4.6.3)
- C++ (g++ 4.6.3)
- PHP (5.3.10)

## API

- `POST /api/evaluate/< language_name >` with file in body.
- `POST /api/evaluate/` with JSON in body. JSON should content
  `language` field & `code` field.

### Curl examples

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

curl  http://192.168.1.2:8080/api/evaluate/cpp -d "
#include <iostream>

int main()
{
std::cout << \"Hello\"; return 0;
}
"
{"stdout":"Hello\n","stderr":"","meta":"","exitCode":0}
```
## Install

```bash
make install
```

Instead of `make install` that just download universal image from docker
images repository, you can build image locally:

```bash
make build-image
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
./bin/perdocker -h
Usage of ./bin/perdocker:
  -c-workers=1: Count of C workers.
  -cpp-workers=1: Count of C++ workers.
  -golang-workers=1: Count of golang workers.
  -listen=":8080": HTTP server bind to address & port. Ex: localhost:80 or :80
  -nodejs-workers=1: Count of nodejs workers.
  -php-workers=1: Count of PHP workers.
  -python-workers=1: Count of python workers.
  -ruby-workers=1: Count of ruby workers.
  -separate=false: Separate workers by languages.
  -timeout=30: Max execution time.
  -workers=1: Count of workers for non separated workers.
```

## Defaults

- listen :8080.
- non separate mode.
- 1 worker per language or 1 worker in non separate mode.
- timeout 30 seconds.

## TODO

- run process inside exist docker container. Required docker 0.8.0.
- restart container after some amount evals.
