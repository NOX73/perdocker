Perdocker
====

Evaluate code in differnet languages inside docker containers.

## Langs

Currently supported languages:

- ruby (2.1.0)

Cooming soon: 

- javascript
- golang

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

