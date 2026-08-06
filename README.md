# secret

[![CI](https://github.com/uchaloop/secret/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/secret/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/secret/v2.svg)](https://pkg.go.dev/github.com/uchaloop/secret/v2)
[![License: MIT](https://img.shields.io/badge/github/license/uchaloop/secret)](LICENSE)

An opaque value for passwords, tokens, and other sensitive strings. It masks
values during formatting, structured logging, and serialization while keeping
access to the original string explicit.

## Installation

```bash
go get github.com/uchaloop/secret/v2
```

```go
import "github.com/uchaloop/secret/v2"
```

## Usage

```go
password := secret.New(raw)

fmt.Println(password)                  // ****
slog.Info("config", "password", password) // password=****

raw = password.Reveal()
```

`Reveal` returns the original string. Keep the returned value scoped to the
component that needs it and do not log or serialize it.

The zero value is ready to use:

```go
var password secret.Secret

password.IsZero() // true
password.String() // ""
```

## Configuration

`Secret` can be populated by environment and text decoders:

```go
type Config struct {
	Password secret.Secret `env:"PASSWORD,required"`
}
```

With `confmaker`, keep secrets out of TOML:

```go
type Config struct {
	Password secret.Secret `koanf:"-" env:"PASSWORD,required"`
}
```

## Masking

A non-empty `Secret` is rendered as `****`; an empty value is rendered as an
empty string. Masking is supported for:

- `fmt` formatting;
- `log/slog`;
- text and JSON serialization;
- XML serialization.

```go
cfg := struct {
	Password secret.Secret `json:"password"`
}{
	Password: secret.New("sensitive"),
}

data, err := json.Marshal(cfg)
// {"password":"****"}
```

Serialization is intentionally lossy and cannot be used to persist a
recoverable secret.

## Check and clear

```go
if password.IsZero() {
	return errors.New("password is required")
}

password.Clear()
```

`Clear` replaces the stored value with an empty string. It cannot guarantee
physical erasure of previous bytes from process memory.

## Sensitive value marker

Use `secret.Value` when an integration needs to recognize sensitive values:

```go
func IsSensitive(value any) bool {
	_, ok := value.(secret.Value)

	return ok
}
```

## Migrating from v1

Update the module path:

```diff
- github.com/uchaloop/secret
+ github.com/uchaloop/secret/v2
```

Use `New` instead of a string conversion:

```diff
- password := secret.Secret(raw)
+ password := secret.New(raw)
```

Call `Reveal` explicitly where the original string is required.

## Security scope

`Secret` prevents accidental disclosure through normal formatting, logging,
and serialization. It does not protect against:

- explicit calls to `Reveal`;
- logging the string returned by `Reveal`;
- `unsafe`, debuggers, crash dumps, or process-memory inspection;
- copies of the original input retained elsewhere.

Treat it as a guardrail, not secure memory.

## License

[MIT](LICENSE)
