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

{"std_out":"1\n4\n9\n","std_err":"","code":0}

curl -POST -d "var a = 6; a += 10; console.log(a)" 'http://localhost:8080/nodejs'
{"std_out":"16\n","std_err":"","code":0}
```

## Install

```bash
make install
```

## Run

```bash
make run
```

## Flags

```bash
make run -port 80 -ruby-workers 5 -nodejs-workers 5
```

## Defaluts

- port 8080
- 1 ruby worker
- 1 nodejs worker

## Coming soon

- timouts per eval.
- many language support (golang or something else).
- improvement run process.
- start & attach to container instead of run it per request.

