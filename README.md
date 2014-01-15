Perdocker
====

Evaluate code in differnet languages inside docker containers.

## Langs

Currently supported languages:

- ruby (2.1.0)

Cooming soon: 

- javascript
- golang

## API

```bash
curl -POST -d "[1,2,3].each do |a| puts a*a; end;" 'http://localhost:8080/ruby'

{"std_out":"1\n4\n9\n","std_err":"","code":0}
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
make run -port 80 -ruby-workers 5
```

## Defaluts

- port 8080
- 1 ruby worker

## Coming soon

- timouts per eval.
- many language support (js, go).
- improvement run process.
- start & attach to container instead of run it per request.

