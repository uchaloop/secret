# secret

[![CI](https://github.com/uchaloop/secret/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/secret/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/secret/v2.svg)](https://pkg.go.dev/github.com/uchaloop/secret/v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`secret` provides an opaque value for passwords, tokens and other sensitive
strings. It protects against accidental disclosure through ordinary formatting,
logging and standard serialization. Access to the underlying string is explicit
through `Reveal`.

The module uses only the Go standard library, so infrastructure libraries can
put `secret.Secret` in their config types without depending on a configuration
or environment-loading stack.

## Install

Version 2 uses the `/v2` module path:

```bash
go get github.com/uchaloop/secret/v2
```

```go
import "github.com/uchaloop/secret/v2"
```

The imported package name remains `secret`.

## Usage

```go
package postgres

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/uchaloop/secret/v2"
)

type Config struct {
	User     string
	Password secret.Secret
}

func Example(raw string) {
	cfg := Config{
		User:     "app",
		Password: secret.New(raw),
	}

	fmt.Printf("%+v\n", cfg)  // {User:app Password:****}
	slog.Info("config", "value", cfg)
	data, _ := json.Marshal(cfg) // {"User":"app","Password":"****"}
	_ = data

	password := cfg.Password.Reveal()
	_ = password // pass it only to the component that needs the real value
}
```

The zero value is ready to use:

```go
var password secret.Secret

password.IsZero() // true
password.String() // ""
```

## Configuration decoding

`Secret` implements `encoding.TextUnmarshaler`, allowing environment and text
decoders to populate it without knowing its internal representation:

```go
type Config struct {
	Password secret.Secret `env:"PASSWORD,required"`
}
```

The tag has no meaning to the `secret` package itself. Source selection and
required-field validation remain the responsibility of the configuration
library.

## API

### Create and access

```go
password := secret.New(raw)
raw = password.Reveal()
```

`Reveal` is intentionally explicit so uses of the real value are easy to find
during review. Do not log or serialize the returned string.

### Check and clear

```go
if password.IsZero() {
	return errors.New("password is required")
}

password.Clear()
```

`Clear` logically replaces the stored value with an empty string. Go strings,
garbage collection and compiler optimizations mean it cannot guarantee physical
erasure of previous bytes from process memory.

### Formatting, logging and serialization

A non-empty `Secret` renders as `****`; an empty one renders as an empty string.
Masking is implemented for:

- `fmt` value verbs through `fmt.Formatter`, `String` and `GoString`;
- `log/slog` through `LogValue`;
- `encoding.TextMarshaler`, including `encoding/json`, through `MarshalText`;
- `encoding/xml` through `MarshalXML`.

Serialization is intentionally lossy: serialized output contains the mask, not
the original value, and should not be used to persist a recoverable secret.

Because `Secret` is an opaque struct rather than a named string, ordinary code
cannot convert, concatenate, index or slice it as a string.

## Sealed marker interface

`Value` identifies sensitive types owned by this module:

```go
func IsSensitive(value any) bool {
	_, ok := value.(secret.Value)
	return ok
}
```

The marker method is unexported, so packages outside `secret` can recognize
implementations but cannot add their own. This lets integration libraries apply
redaction or source policies without depending on the representation of
`Secret`.

## Migrating from v1

Update the module and import path:

```diff
- github.com/uchaloop/secret
+ github.com/uchaloop/secret/v2
```

Replace string-style construction:

```diff
- password := secret.Secret(raw)
+ password := secret.New(raw)
```

`Reveal` and `UnmarshalText` remain available. Code that converted, concatenated,
indexed or sliced `Secret` must now call `Reveal` explicitly at the point where
the real string is required.

## Security scope

`Secret` is designed to prevent accidental disclosure through normal application
code. It does not protect against:

- an explicit call to `Reveal`;
- logging or serializing the string returned by `Reveal`;
- `unsafe` or deliberate reflection-based inspection;
- debuggers, crash dumps or process-memory inspection;
- copies of the original input string retained elsewhere.

Treat `Secret` as a guardrail, not as secure memory.

## Testing

```bash
go test ./...
go test -race ./...
go vet ./...
```

## License

[MIT](LICENSE).
