# lexigen

Console utility that prints random English phrases made of one adjective and one noun.

## Usage

```text
lexigen          # one phrase (default)
lexigen 25       # 25 phrases, one per line on stdout
```

If `count` is greater than the number of words in the smaller dictionary table, the program exits with an error and prints nothing to stdout.

## Build

Requires [Go](https://go.dev/) 1.22 or newer.

Go has no built-in “Gradle”; this repo uses a **Makefile** and small **build scripts** with the same targets.

**One command — Windows + Linux (amd64)** into `dist/`:

```bash
make all
# or (Git Bash / WSL / Linux / macOS):
./scripts/build.sh all
# or PowerShell:
.\scripts\build.ps1 all
```

Other targets: `local`, `windows`, `linux`, `linux-arm64`, `test`, `clean`. Run `make help` or `./scripts/build.sh help`.

Manual single-platform build:

```bash
go build -o lexigen .
```

After a release tag is published:

```bash
go install github.com/Snowbridge/lexigen@latest
```

The binary is placed in `$GOPATH/bin` or `$HOME/go/bin` (ensure that directory is on your `PATH`).

## Install (manual)

Copy the built `lexigen` executable into any directory listed in your `PATH`.

## Data

Word lists are embedded from `words.sqlite` at build time (`adjectives` and `nouns` tables, `word` column).
