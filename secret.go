/*
Package secret provides an opaque, masked value for passwords, tokens and other
sensitive strings. A Secret protects against accidental disclosure through
ordinary formatting, logging and serialization; the real value leaves the type
only through an explicit Reveal.

	password := secret.New(raw)

	fmt.Println(password)                      // ****
	slog.Info("config", "password", password)  // password=****

	raw = password.Reveal()

The zero value is an empty secret, ready to use: IsZero reports true and every
representation of it is the empty string rather than a mask, so an unset value
does not read as a set one.

# What is masked

A non-empty Secret renders as **** everywhere a value is rendered without being
asked for: fmt verbs, log/slog, text and JSON, XML. An empty one renders as the
empty string.

	cfg := struct {
		Password secret.Secret `json:"password"`
	}{Password: secret.New("sensitive")}

	data, _ := json.Marshal(cfg)  // {"password":"****"}

Serialization is lossy on purpose. A Secret cannot be round-tripped through
JSON: what comes out is the mask, so a marshalled config cannot be used to carry
the value on. UnmarshalText, on the other hand, reads a value in - it is what
lets an environment or text decoder populate a Secret while every outward
representation stays masked.

Format leaves the type and pointer verbs alone. %T and %p say what fmt already
knows without going near the contents.

# Reading it back

Reveal returns the underlying string, and it is the only way out. That is the
point: every real use of a secret is one greppable call, so an audit is a search
rather than a reading. Keep what it returns scoped to the component that needs
it, and neither log nor serialize it.

Clear empties a Secret logically. It cannot guarantee that the bytes of the
previous value are gone from process memory - Go strings are immutable and may
have been copied - so treat it as intent, not erasure.

# Recognising a secret

Value is a sealed marker: the types this package owns implement it, and no other
package can, because the marker method is unexported. An integration library
consults it to treat a field differently without importing a notion of secrecy
of its own:

	func IsSensitive(value any) bool {
		_, ok := value.(secret.Value)

		return ok
	}

A config dumper does this to report a field as set or unset instead of printing
what it holds.

# In a config

The package has no dependencies outside the standard library, so an
infrastructure library can put a Secret in a config struct without taking on a
configuration or environment-parsing stack:

	type Config struct {
		Password secret.Secret `env:"PASSWORD,notEmpty"`
	}

# What it does not protect against

A Secret is a guardrail against disclosure by accident, not secure memory. It
does nothing about an explicit Reveal, about logging the string Reveal returned,
about unsafe, a debugger, a crash dump or anything else reading process memory,
or about a copy of the original input kept somewhere else.
*/
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
