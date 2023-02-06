# Spider

Spider is a WordPress website (blog) deployment tool. It Automates the process of launching and transfering database assets to a new WordPress blog site.

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
[build_location]/[program_name] [flag] [new_website_slug]
```

## Example

### Staging to Production:

```console
~/Documents/programs/spider -s2p antiracism
```

## License

Code is distributed under [The Unlicense](https://github.com/nausicaan/spider/blob/main/LICENSE.md) and is part of the Public Domain.
