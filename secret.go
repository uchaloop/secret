// Package secret provides a masked string type for passwords, tokens and other
// sensitive values. A struct or config that carries a Secret never leaks it when
// formatted, logged, or serialized; the real value is obtained only through an
// explicit Reveal call.
//
// The package has no external dependencies, so infrastructure libraries can use
// Secret in their (source-agnostic) Config structs without pulling in a config
// or env-parsing dependency.
package secret

import "log/slog"

// mask is the placeholder shown in every textual representation of a non-empty
// Secret.
const mask = "****"

// Secret is a sensitive string that masks itself everywhere it could be exposed:
// fmt verbs (%v/%s/%+v/%#v), encoding/json and encoding/text (via MarshalText),
// and log/slog (via LogValue). An empty Secret renders as empty.
//
// The underlying value is retrieved only with Reveal, which makes every real use
// greppable.
type Secret string

// Reveal returns the underlying secret value. Call it only where the real value
// is genuinely required (e.g. building a connection); never log the result.
func (s Secret) Reveal() string {
	return string(s)
}

// String masks the value for fmt %v/%s and any fmt.Stringer consumer.
func (s Secret) String() string {
	if len(s) == 0 {
		return ""
	}

	return mask
}

// GoString masks the value for fmt %#v.
func (s Secret) GoString() string {
	if len(s) == 0 {
		return `""`
	}

	return `"` + mask + `"`
}

// MarshalText masks the value for encoding/json, TOML and any other
// encoding.TextMarshaler-based encoder, so serializing a struct never leaks it.
func (s Secret) MarshalText() ([]byte, error) {
	if len(s) == 0 {
		return []byte{}, nil
	}

	return []byte(mask), nil
}

// LogValue masks the value for log/slog structured logging.
func (s Secret) LogValue() slog.Value {
	if len(s) == 0 {
		return slog.StringValue("")
	}

	return slog.StringValue(mask)
}

// UnmarshalText fills the Secret with the raw bytes, so an env or text decoder
// (for example a config loader reading a value out of the environment) can set
// it. It is deliberately asymmetric with MarshalText: the real value is read in
// here, while every outward representation stays masked.
func (s *Secret) UnmarshalText(text []byte) error {
	*s = Secret(text)

	return nil
}
