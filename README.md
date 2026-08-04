# secret

[![CI](https://github.com/uchaloop/secret/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/secret/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/secret.svg)](https://pkg.go.dev/github.com/uchaloop/secret)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A masked string type for passwords, tokens and other sensitive values. A struct
or config that carries a `Secret` never leaks it when formatted, logged, or
serialized; the real value is obtained only through an explicit `Reveal` call.

**Zero dependencies** (standard library only), so any library can hold a `Secret`
in its config struct without pulling in a config or env-parsing stack.

## Install

```bash
go get github.com/uchaloop/secret
```

## Usage

```go
type Config struct {
	User     string
	Password secret.Secret
}

cfg := Config{User: "app", Password: secret.Secret(raw)}

log.Info("db", "config", cfg)        // password logs as "****"
fmt.Sprintf("%+v", cfg)              // "****"
json.Marshal(cfg)                    // "****"

dsn := cfg.Password.Reveal()         // the real value, only at the point of use
```

`Secret` masks itself in:

- `fmt` verbs (`%v`/`%s`/`%q`/`%+v`/`%#v`/`%x`) via `String`/`GoString`,
- `encoding/json` and `encoding/text` via `MarshalText`,
- `log/slog` via `LogValue`.

An empty `Secret` renders as empty. `UnmarshalText` reads a value in (for an env
or text decoder) while every outward representation stays masked. The real value
leaves the type only through `Reveal`, which makes every real use greppable.

## Testing

```bash
go test ./...
```

## License

[MIT](LICENSE).
