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
./wp2go -a example-szkolaplywania.pl.tar.gz -d baza-2025-02-05.sql -s example
```

### Options

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--archive` | `-a` | yes | Path to WordPress .tar.gz archive |
| `--db` | `-d` | yes | Path to SQL database dump |
| `--sitename` | `-s` | yes | Site name (e.g. `example` → https://example.ddev.site) |
| `--output-path` | `-o` | no | Directory where to create the project (default: current directory) |
| `--ui` | | no | Use interactive prompts for missing options |
| `--force` | `-f` | no | Overwrite existing project directory |
| `--config` | | no | Path to config file (yaml) |
| `--host` | | no | Host name in config (under `hosts.<name>`) |

### Examples

Create project in current directory (project dir will be `./example/`):

```bash
wp2go -a archive.tar.gz -d dump.sql -s example
```

Create project under `~/src/example`:

```bash
wp2go -a archive.tar.gz -d dump.sql -s example -o ~/src/example
```

Interactive mode (prompts for archive, db, sitename, output path if not provided):

```bash
wp2go --ui
```

Force overwrite existing project:

```bash
wp2go -a archive.tar.gz -d dump.sql -s delfinki --force
```

### Config file

You can define host-specific settings in YAML and select them with `--host`.
Config files are searched in `./wp2go.{yaml,yml}` and `~/.config/wp2go/wp2go.{yaml,yml}`
unless `--config` is provided.
If a host is selected and `sitename` is omitted, it defaults to the host name.
You can also set `host: <name>` in the config to choose a default host.

Example:

```yaml
defaults:
  output_path: ~/src
  force: false

host: delfinki

hosts:
  delfinki:
    archive: /backups/delfinki.tar.gz
    db: /backups/delfinki.sql
    sitename: delfinki

  example:
    archive: /backups/example.tar.gz
    db: /backups/example.sql
    sitename: example
    output_path: ~/src/example
    force: true
```

Usage:

```bash
wp2go --host delfinki
```

```bash
wp2go --config ~/wp2go.yml --host delfinki
```

## What it does

1. Creates a directory named after the sitename (e.g. `example`) in the output path (or current dir).
2. Extracts the .tar.gz archive there; if the archive has a single top-level directory, its contents are placed at project root.
3. Patches `wp-config.php` so it uses DDEV’s database credentials (`DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_HOST` → `db`); `$table_prefix` is left unchanged.
4. Runs `ddev config --project-type=wordpress --project-name=<sitename>` and `ddev start`.
5. Detects the old site URL from the SQL dump (from `*_options` `siteurl`/`home`), replaces it with the DDEV URL, and runs `ddev import-db`.

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
│   ├── sql/          # Detect old URL, replace, temp file for import
│   └── wpconfig/     # Patch wp-config.php for DDEV DB credentials
├── main.go
├── go.mod
└── go.sum
```
