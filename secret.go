// Package secret provides an opaque, masked value for passwords, tokens and
// other sensitive strings. Secret protects against accidental disclosure
// through ordinary formatting, logging and standard serialization; the real
// value is obtained only through an explicit Reveal call.
//
// The package has no external dependencies, so infrastructure libraries can use
// Secret in source-agnostic config structs without pulling in a configuration or
// environment-parsing stack.
package secret

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"strconv"
)

// mask is the placeholder shown for every non-empty Secret.
const mask = "****"

// Value is a sealed marker for sensitive value types owned by this package.
// Other packages can recognize these values, but cannot add implementations
// because the marker method is intentionally unexported.
type Value interface {
	isSecret()
}

// Secret is an opaque sensitive string. Its zero value is an empty secret.
//
// Secret masks itself during formatting, logging and text serialization. Its
// contents cannot be accessed through ordinary string operations; use Reveal
// explicitly at the point where the underlying value is genuinely required.
type Secret struct {
	value string
}

var _ Value = Secret{}

// New returns a Secret containing value.
func New(value string) Secret {
	return Secret{value: value}
}

func (Secret) isSecret() {}

// Reveal returns the underlying value. Call it only where the real secret is
// required, and never log or serialize the returned string.
func (s Secret) Reveal() string {
	return s.value
}

// IsZero reports whether s is empty.
func (s Secret) IsZero() bool {
	return len(s.value) == 0
}

// Clear logically empties s. It does not guarantee physical erasure of prior
// string data from process memory.
func (s *Secret) Clear() {
	if s != nil {
		s.value = ""
	}
}

// String returns an empty string for an empty Secret and a mask otherwise.
func (s Secret) String() string {
	return s.masked()
}

// GoString returns a masked Go-syntax representation.
func (s Secret) GoString() string {
	return strconv.Quote(s.masked())
}

// Format masks Secret for fmt value verbs. Type and pointer verbs remain safe:
// fmt handles %T and pointer %p itself without exposing the contents.
func (s Secret) Format(state fmt.State, verb rune) {
	value := s.masked()
	if verb == 'q' {
		value = strconv.Quote(value)
	}

	_, _ = io.WriteString(state, value)
}

// MarshalText returns an empty value for an empty Secret and a mask otherwise.
// Serializing a Secret is intentionally lossy and never exposes its contents.
func (s Secret) MarshalText() ([]byte, error) {
	return []byte(s.masked()), nil
}

// MarshalXML writes only the masked representation.
func (s Secret) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(s.masked(), start)
}

// LogValue returns a masked slog value.
func (s Secret) LogValue() slog.Value {
	return slog.StringValue(s.masked())
}

// UnmarshalText replaces s with the supplied input. It enables environment and
// text decoders to populate a Secret while all outward representations remain
// masked.
func (s *Secret) UnmarshalText(text []byte) error {
	s.value = string(text)

	return nil
}

func (s Secret) masked() string {
	if s.IsZero() {
		return ""
	}

	return mask
}
