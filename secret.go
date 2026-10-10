/*
Package secret provides an opaque, masked value for passwords, tokens and other
sensitive strings. A Secret protects against accidental disclosure through
supported formatting, logging and serialization. Reveal is the explicit API
for reading the underlying value.

	password := secret.New(raw)

	fmt.Println(password)                      // ****
	slog.Info("config", "password", password)  // password=****

	raw = password.Reveal()

The zero value is an empty secret, ready to use. IsZero reports true and its
formatting and serialization methods return an empty representation.

# What is masked

A non-empty Secret renders as **** through its fmt, log/slog, text, JSON and
XML interfaces. An empty one renders as the empty string. When fmt traverses
unexported fields, it may bypass these methods and show a technical representation
instead of the mask. The stored secret is not disclosed in that representation.

	cfg := struct {
		Password secret.Secret `json:"password"`
	}{Password: secret.New("sensitive")}

	data, _ := json.Marshal(cfg)  // {"password":"****"}

Serialization is lossy on purpose. A Secret cannot be round-tripped through
JSON: what comes out is the mask, so a marshalled config cannot be used to carry
the value on. UnmarshalText, on the other hand, reads a value in - it is what
lets an environment or text decoder populate a Secret while outward serialization
stays masked. UnmarshalJSON accepts a string or null; null clears the receiver.
Other JSON values, invalid UTF-8 and unpaired UTF-16 surrogate escapes are rejected.
Nil receivers return errors from both input methods.
Errors from UnmarshalJSON do not expose input or decoder causes and leave the
receiver unchanged.

An enclosing decoder can fail before or after calling UnmarshalJSON, and its
errors are outside this guarantee. For example, a streaming decoder can assign
one value before rejecting trailing data. To preserve an existing secret on any
document error, decode into a temporary Secret and assign it only after success.

Format leaves the type and pointer verbs alone. %T and %p say what fmt already
knows without going near the contents.

# Reading it back

Reveal returns the underlying string. Use it explicitly at the point that needs
the value, so ordinary accesses can be found during review. Neither log nor
serialize the returned string. Reflection can bypass this API; the type is not
a boundary against code deliberately inspecting its representation.

Clear empties a Secret logically. It cannot guarantee that the bytes of the
previous value are gone from process memory - Go strings are immutable and may
have been copied - so treat it as intent, not erasure. Clearing or replacing one
Secret does not change previously made copies. Concurrent reads are supported;
mutation of the same variable through Clear, UnmarshalText or UnmarshalJSON
must be synchronized with other accesses by the caller.

# Recognising a secret

Secret implements interface{ IsSensitive() }, allowing integrations to recognize
sensitivity without importing this package. The method is an empty marker, not a
predicate, and does not inspect the value. Integrations may inspect the type
without calling it.

Value remains supported as a marker with an unexported method. Other packages cannot implement
that method directly, but can inherit it by embedding Secret or Value. It is a
sensitivity hint for integrations, not proof of a type's origin or safe formatting:

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
about reflection, unsafe, a debugger, a crash dump or anything else reading
process memory, or about a copy of the original input kept somewhere else.
*/
package secret

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"unicode/utf8"
)

const mask = "****"

// Value marks sensitive values for integrations. Its unexported method cannot
// be implemented directly by another package, but can be inherited by embedding.
// The marker does not guarantee safe formatting by an embedding type.
type Value interface {
	isSecret()
}

// Secret is an opaque sensitive string. Its zero value is an empty secret.
//
// Secret masks itself during formatting, logging and text serialization. Its
// contents cannot be accessed through ordinary string operations; use Reveal
// explicitly at the point where the underlying value is required.
// Secret is not comparable: use IsZero to check emptiness. It cannot be a map key
// or compared with ==, including against its zero value.
type Secret struct {
	reveal func() string
}

var _ Value = Secret{}

// New returns a Secret containing value.
func New(value string) Secret {
	if len(value) == 0 {
		return Secret{}
	}

	// fmt bypasses methods on private fields. A closure keeps captured text out
	// of its reflective output, including diagnostics for unsupported verbs.
	return Secret{reveal: func() string { return value }}
}

func (Secret) isSecret() {}

// IsSensitive marks Secret for integrations without a package dependency.
// It does not reveal or change the value.
func (Secret) IsSensitive() {}

// Reveal returns the underlying value. Call it only where the real secret is
// required, and never log or serialize the returned string.
func (s Secret) Reveal() string {
	if s.reveal == nil {
		return ""
	}

	return s.reveal()
}

// IsZero reports whether s is empty.
func (s Secret) IsZero() bool {
	return s.reveal == nil
}

// Clear logically empties s. It does not guarantee physical erasure of prior
// string data from process memory. Existing copies retain their values.
func (s *Secret) Clear() {
	if s == nil {
		return
	}

	s.reveal = nil
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
	maskedValue := s.masked()
	if verb == 'q' {
		maskedValue = strconv.Quote(maskedValue)
	}

	_, _ = io.WriteString(state, maskedValue)
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
// masked. Input bytes are preserved without Unicode validation and copied, so
// subsequent changes to text do not affect s. A nil receiver returns an error.
func (s *Secret) UnmarshalText(text []byte) error {
	if s == nil {
		return errors.New("secret: cannot decode into a nil receiver")
	}

	*s = New(string(text))

	return nil
}

// UnmarshalJSON replaces s with a JSON string, or clears it for null.
// Invalid input, including malformed Unicode, and nil receivers return errors.
// On error s is unchanged.
// Errors returned by this method contain neither input nor decoder causes.
// An enclosing decoder may fail before or after calling this method; decode into
// a temporary value for atomic replacement on document-level success.
func (s *Secret) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("secret: cannot decode into a nil receiver")
	}

	if bytes.Equal(bytes.Trim(data, " \t\r\n"), []byte("null")) {
		s.Clear()

		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return errors.New("secret: expected a JSON string or null")
	}

	if !validJSONUnicode(data) {
		return errors.New("secret: expected a JSON string or null")
	}

	*s = New(value)

	return nil
}

// validJSONUnicode checks a syntactically valid JSON string before accepting
// encoding/json's replacement of malformed Unicode.
func validJSONUnicode(data []byte) bool {
	if !utf8.Valid(data) {
		return false
	}

	for index := 0; index < len(data); index++ {
		if data[index] != '\\' {
			continue
		}

		index++
		if data[index] != 'u' {
			continue
		}

		codeUnit, err := strconv.ParseUint(string(data[index+1:index+5]), 16, 16)
		if err != nil {
			return false
		}

		index += 4
		if codeUnit < 0xd800 || codeUnit > 0xdfff {
			continue
		}

		if codeUnit > 0xdbff || len(data)-index <= 6 || data[index+1] != '\\' || data[index+2] != 'u' {
			return false
		}

		lowSurrogate, err := strconv.ParseUint(string(data[index+3:index+7]), 16, 16)
		if err != nil {
			return false
		}

		if lowSurrogate < 0xdc00 || lowSurrogate > 0xdfff {
			return false
		}

		index += 6
	}

	return true
}

func (s Secret) masked() string {
	if s.IsZero() {
		return ""
	}

	return mask
}
