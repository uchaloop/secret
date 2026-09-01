# secret

[![CI](https://github.com/uchaloop/secret/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/secret/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/secret/v2.svg)](https://pkg.go.dev/github.com/uchaloop/secret/v2)
[![License: MIT](https://img.shields.io/github/license/uchaloop/secret)](LICENSE)

An opaque value for passwords, tokens and other sensitive strings. It masks
itself wherever a value is rendered without being asked for, and the real string
leaves the type only through an explicit call.

- **One way out** - `Reveal` is the only one, so every real use of a secret is a
  greppable call and an audit is a search.
- **Masked by default** - `fmt`, `log/slog`, text, JSON and XML, all of it.
- **Opaque** - a struct, not a string: no conversion, concatenation, indexing or
  slicing.
- **No dependencies**, so an infrastructure library can put one in a config
  struct without taking on a configuration stack.

```bash
go get github.com/uchaloop/secret/v2
```

## Quick start

```go
password := secret.New(raw)

fmt.Println(password)                      // ****
slog.Info("config", "password", password)  // password=****

raw = password.Reveal()
```

In a config a library declares and never reads itself:

```go
type Config struct {
	Password secret.Secret `env:"PASSWORD,notEmpty"`
}
```

## What it is not

A guardrail against disclosure by accident, not secure memory. It does nothing
about an explicit `Reveal`, about logging what `Reveal` returned, or about
`unsafe`, a debugger, a crash dump or anything else that reads process memory.

## Documentation

What is masked, what `Reveal` and `Clear` promise, and the reasons behind them
are in the package documentation:
**[pkg.go.dev/github.com/uchaloop/secret/v2](https://pkg.go.dev/github.com/uchaloop/secret/v2)**.

Upgrading from v1 is in the [changelog](CHANGELOG.md).

## License

[MIT](LICENSE)
