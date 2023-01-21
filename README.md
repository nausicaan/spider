# Deploy WordPress Site

Automate the process of launching a site and transfering database assets to a new WordPress Domain.

## Prerequisites

- Googles' [Go language](https://go.dev) installed to enable building executables from source code.

## Build

From the root folder containing *main.go*, use the command that matches your environment:

### Windows & Mac:

```bash
go build -o <build_location>/<program_name> main.go
```

### Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o <build_location>/<program_name> main.go
```

## Run

```bash
<build_location>/<program_name> <flag> <new_website_slug>
```

## Example

### Staging to Production:

```bash
~/Documents/programs/deploy -s2p antiracism
```

## License

Code is distributed under [The Unlicense](https://github.com/nausicaan/deploy/blob/main/LICENSE.md) and is part of the Public Domain.
