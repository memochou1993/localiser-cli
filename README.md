# Localiser CLI

Localiser CLI is a command-line tool that lets you sync locale files from server to local.

## Usage

Download [binary](https://github.com/memochou1993/localiser-cli/tree/master/bin) and set to `PATH`.

Create `localiser.yaml` to specified project.

```YAML
---
endpoint: http://localhost:8000/api
project_id: 1
output_directory: src/assets/lang
```

Sync locale files with `localiser` command.

```BASH
localiser
```

## Development

Download the master branch.

```BASH
git clone git@github.com:memochou1993/localiser-cli.git
```

Copy `localiser.example.yaml` to `localiser.yaml`.

```BASH
cp localiser.example.yaml localiser.yaml
```

Download locale files from server.

```BASH
go run main.go
```
