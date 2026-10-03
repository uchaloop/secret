# secret

[![CI](https://github.com/uchaloop/secret/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/secret/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/gh/uchaloop/secret/branch/main/graph/badge.svg)](https://app.codecov.io/gh/uchaloop/secret)
[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/secret/v2.svg)](https://pkg.go.dev/github.com/uchaloop/secret/v2)
[![License: MIT](https://img.shields.io/github/license/uchaloop/secret)](LICENSE)

An opaque value for passwords, tokens and other sensitive strings. It masks
itself through supported formatting and serialization interfaces and provides
explicit access to the underlying string.

Requires Go 1.21.0 or later for `log/slog` integration. Use a currently supported
Go toolchain for application builds; this minimum describes source compatibility.

- **Explicit access** - `Reveal` makes ordinary uses of the underlying string
  easy to find during review.
- **Masked by default** - through `fmt`, `log/slog`, text, JSON and XML interfaces.
  Private fields may display a technical representation instead of a mask,
  without disclosing the stored value.
- **Opaque** - a struct, not a string: no conversion, concatenation, indexing or
  slicing.
- **No dependencies**, so an infrastructure library can put one in a config
  struct without taking on a configuration stack.

```bash
go get github.com/uchaloop/secret/v2
```

## Quick start

```go
import "github.com/uchaloop/secret/v2"
```

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

`Secret` is not comparable: use `IsZero()` rather than comparing against
`Secret{}`. Comparing secrets with `==` or using them as map keys is unsupported.
Comparing interface values containing secrets can panic, as with other
non-comparable Go values. No equality method is provided.

## JSON input

JSON strings replace the secret; `null` clears it. Other values, invalid UTF-8
and unpaired UTF-16 surrogate escapes are rejected instead of silently changing
the secret. Only JSON whitespace (space, tab, CR and LF) is accepted.
`UnmarshalText` preserves arbitrary input bytes without Unicode validation.
`UnmarshalJSON` errors do not include input or decoder causes and leave the
receiver unchanged. Both `UnmarshalJSON` and `UnmarshalText` reject nil receivers.

An outer decoder can fail after assigning a value, for example on trailing data.
For atomic replacement of an existing secret, decode into a temporary value:

```go
var next secret.Secret
if err := json.Unmarshal(data, &next); err != nil {
	return err
}

password = next
```

This pattern also applies to `encoding/json/v2`. Errors from an outer decoder
may contain input details; do not log them as though they were sanitized by secret.
Serialization remains intentionally lossy: it writes a mask, never the secret.

## What it is not

A guardrail against disclosure by accident, not secure memory. It does nothing
about an explicit `Reveal`, about logging what `Reveal` returned, or about
reflection, `unsafe`, a debugger, a crash dump or anything else that reads process memory.
`Clear` affects only its receiver; existing copies retain their values. Synchronize
mutation of a shared variable with other accesses.

## Documentation

What is masked, what `Reveal` and `Clear` promise, and the reasons behind them
are in the package documentation:
**[pkg.go.dev/github.com/uchaloop/secret/v2](https://pkg.go.dev/github.com/uchaloop/secret/v2)**.

## Upgrading to 2.1.0

This release keeps the `/v2` import path but includes incompatible behavior changes:

- Replace comparisons with `Secret{}` by `IsZero()`. `Secret` is no longer
  comparable and cannot be a map key; comparisons through interfaces can panic.
- JSON `null` now clears a secret. Decode into a temporary value if the entire
  document must succeed before replacing an existing value.
- Malformed Unicode in JSON is rejected rather than silently replaced. Text
  input remains byte-preserving.

See the [changelog](CHANGELOG.md) for release details and the earlier v1 migration.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for project scope, development checks,
code style and pull request conventions. For suspected disclosures, follow
[SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
