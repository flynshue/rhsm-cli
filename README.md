# rhsm-cli

# Getting a token
Please see [Generating an offline token](https://access.redhat.com/management/api) to generate a token.

```bash
export RHSM_TOKEN=<your offline token>
```

# Building
I'll set up goreleaser later, but for now just do a go build
```bash
go build .
```

# Using rhsm-cli

## Listing all systems under account
```bash
./rhsm-cli list systems
```

## List systems matching a filter
```bash
./rhsm-cli list systems --filter ocp
```