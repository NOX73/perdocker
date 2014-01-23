Perdocker
====

Evaluate code in differnet languages inside docker containers.

## Langs

Currently supported languages:

- ruby (2.1.0)
- javascript (nodejs 0.10.24)

Cooming soon: 

- golang

## API

```bash
curl -POST -d "[1,2,3].each do |a| puts a*a; end;" 'http://localhost:8080/ruby'
{"stdout":"1\n4\n9\n","stderr":"","exitCode":0}

curl -POST -d "var a = 6; a += 10; console.log(a)" 'http://localhost:8080/nodejs'
{"stdout":"16\n","stderr":"","exitCode":0}
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
> Perdocker correctly works only with latest dev Docker version. Cause
> is [this bug](https://github.com/dotcloud/docker/issues/1319). Bug
> will be fixed only for 0.8.0 version of Docker. But compiled dev vesion
> of docker you can find in `bin/` directory.
> Yout can install latest docker version and then just replace
> `/usr/bin/docker` with `./bin/docker` and restart service.

> **NOTE:**
> Perdocker expects that it can run `docker` command without sudo.
> [For details](http://docs.docker.io/en/latest/use/basics/)

## Flags

```bash
./bin/perdocker -port 80 -ruby-workers 5 -nodejs-workers 5 -timeout 5
```

## Defaluts

- port 8080
- 1 ruby worker
- 1 nodejs worker
- timeout 60 seconds

## Coming soon

- timouts per eval.
- many language support (golang, php, C, C++, python and something else).
- improvement run process.
- start & attach to container instead of run it per request.

