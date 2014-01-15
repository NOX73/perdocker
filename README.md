Perdocker
====

Evaluate code in differnet languages inside docker containers.

## Install

```bash
make images_pull
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
