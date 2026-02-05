# wp2go

Restore a WordPress site from a `.tar.gz` archive and SQL dump into a DDEV project.

## Requirements

- Go 1.21+
- [DDEV](https://ddev.com/) installed and on `PATH`

## Build

```bash
go mod download
go build -o wp2go .
```

## Usage

```bash
./wp2go -a delfinki-szkolaplywania.pl.tar.gz -d baza-2025-02-05.sql -s delfinki
```

### Options

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--archive` | `-a` | yes | Path to WordPress .tar.gz archive |
| `--db` | `-d` | yes | Path to SQL database dump |
| `--sitename` | `-s` | yes | Site name (e.g. `delfinki` → https://delfinki.ddev.site) |
| `--output-path` | `-o` | no | Directory where to create the project (default: current directory) |
| `--ui` | | no | Use interactive prompts for missing options |

### Examples

Create project in current directory (project dir will be `./delfinki/`):

```bash
wp2go -a archive.tar.gz -d dump.sql -s delfinki
```

Create project under `~/src/delfinki`:

```bash
wp2go -a archive.tar.gz -d dump.sql -s delfinki -o ~/src/delfinki
```

Interactive mode (prompts for archive, db, sitename, output path if not provided):

```bash
wp2go --ui
```

## What it does

1. Creates a directory named after the sitename (e.g. `delfinki`) in the output path (or current dir).
2. Extracts the .tar.gz archive there; if the archive has a single top-level directory, its contents are placed at project root.
3. Runs `ddev config --project-type=wordpress --project-name=<sitename>` and `ddev start`.
4. Detects the old site URL from the SQL dump (from `*_options` `siteurl`/`home`), replaces it with the DDEV URL, and runs `ddev import-db`.

## Project layout

```
wp2go/
├── cmd/
│   └── root.go       # Cobra root command, flags, promptui when --ui
├── internal/
│   ├── archive/      # Extract .tar.gz, flatten single root dir
│   ├── ddev/         # ddev config, start, import-db
│   ├── domain/       # Sitename → short name + FQDN
│   ├── restore/      # Orchestration
│   └── sql/          # Detect old URL, replace, temp file for import
├── main.go
├── go.mod
└── go.sum
```
