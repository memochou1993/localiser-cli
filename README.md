# Localiser CLI

Localiser CLI is a command-line tool that downloads locale files from Localiser to local.

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

```
go run main.go
```
