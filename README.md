# Spider

Spider is a WordPress website (blog) deployment tool. It Automates the process of launching and transfering database assets to a new WordPress blog site.

![Spider](spider.webp)

*Image by [relato](https://www.vectorstock.com/royalty-free-vectors/vectors-by_relato) on [VectorStock](https://www.vectorstock.com)*

## Prerequisites

- Googles' [Go language](https://go.dev) installed to enable building executables from source code.

## Build

From the root folder containing *main.go*, use the command that matches your environment:

### Windows & Mac:

```bash
go build -o [name] main.go
```

### Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o [name] main.go
```

## Run

```bash
[program] [flag] [website_slug]
```

## Flags

Current flages are:

- s2p - Staging to Production and

- p2s - Production to Staging

Example deployment:

```bash
spider -s2p antiracism
```

## License

Code is distributed under [The Unlicense](https://github.com/nausicaan/spider/blob/main/LICENSE.md) and is part of the Public Domain.
