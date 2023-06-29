# Spider

Spider is a WordPress website (blog) deployment tool. It Automates the process of launching and transfering database assets to a new WordPress blog site.

![Spider](spider.jpg)

## Prerequisites

- Googles' [Go language](https://go.dev) installed to enable building executables from source code.

## Build

From the root folder containing *main.go*, use the command that matches your environment:

### Windows & Mac:

```console
go build -o [name] main.go
```

### Linux:

```console
GOOS=linux GOARCH=amd64 go build -o [name] main.go
```

## Run

```console
[program] [flag] [website_slug]
```

## Flags

Current flages are:

- s2p - Staging to Production and

- p2s - Production to Staging

Example deployment:

```console
spider -s2p antiracism
```

## License

Code is distributed under [The Unlicense](https://github.com/nausicaan/spider/blob/main/LICENSE.md) and is part of the Public Domain.
