[![Build](https://github.com/stanislavbebej-ext43345/summon-secrets-manager/actions/workflows/build.yml/badge.svg)](.github/workflows/build.yml)
[![dependabot](https://img.shields.io/badge/Dependabot-enabled-brightgreen?logo=dependabot)](.github/dependabot.yml)
[![release-please](https://img.shields.io/badge/release--please-enabled-brightgreen?logo=google)](release-please-config.json)

# summon-secrets-manager

[Bitwarden Secrets Manager](https://bitwarden.com/products/secrets-manager/) provider for [Summon](https://github.com/cyberark/summon) in [go](https://github.com/bitwarden/sdk-sm/tree/main/languages/go).

## Development

```bash
export BINARY_NAME="summon-secrets-manager"

go env -w CGO_ENABLED=1
go env -w CC=musl-gcc

go build -ldflags '-linkmode external -extldflags "-static -Wl,-unresolved-symbols=ignore-all"' -o $BINARY_NAME
strip $BINARY_NAME
upx -q -9 $BINARY_NAME
sudo cp $BINARY_NAME /usr/local/bin/Providers
```

## Usage

1. create a [secrets.yml](./secrets.yml) configuration file with `secretId`s:

```yaml
# Retrieve the "password" field
# This is equivalent to "acd2d25f-1fd2-4604-9a6e-b2a600f71a31:Value"
SECRET_VARIABLE: !var acd2d25f-1fd2-4604-9a6e-b2a600f71a31

# Retrieve other field, e.g.: the "key"
USERNAME_VARIABLE: !var acd2d25f-1fd2-4604-9a6e-b2a600f71a31:Key
```

2. run `summon`:

```bash
export BWS_ACCESS_TOKEN='***'

summon -p summon-keepass printenv
```

Supported fields (case-insensitive):

- `Key`
- `Note`
- `Value`
