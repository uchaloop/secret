package secret

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

const value = "hunter2"

func TestNewRevealAndZero(t *testing.T) {
	var zero Secret
	if !zero.IsZero() || zero.Reveal() != "" {
		t.Fatal("zero Secret must be empty")
	}

	s := New(value)
	if s.IsZero() {
		t.Fatal("New returned an empty Secret")
	}
	if got := s.Reveal(); got != value {
		t.Fatalf("Reveal() = %q, want original value", got)
	}
}

func TestClear(t *testing.T) {
	s := New(value)
	s.Clear()
	if !s.IsZero() || s.Reveal() != "" {
		t.Fatal("Clear did not empty Secret")
	}

	var nilSecret *Secret
	nilSecret.Clear()
}

func TestImplementsSealedMarker(t *testing.T) {
	var _ Value = Secret{}
	var _ Value = (*Secret)(nil)
}

func TestMasksAllFmtVerbs(t *testing.T) {
	s := New(value)
	forms := map[string]string{
		"%v":     fmt.Sprintf("%v", s),
		"%s":     fmt.Sprintf("%s", s),
		"%q":     fmt.Sprintf("%q", s),
		"%+v":    fmt.Sprintf("%+v", s),
		"%#v":    fmt.Sprintf("%#v", s),
		"%x":     fmt.Sprintf("%x", s),
		"%X":     fmt.Sprintf("%X", s),
		"%d":     fmt.Sprintf("%d", s),
		"append": string(fmt.Appendf(nil, "%+v", s)),
	}
	for verb, out := range forms {
		if strings.Contains(out, value) {
			t.Errorf("%s leaks the secret: %q", verb, out)
		}
		if !strings.Contains(out, mask) {
			t.Errorf("%s = %q, want masked output", verb, out)
		}
	}

	if out := fmt.Sprintf("%p", &s); strings.Contains(out, value) {
		t.Fatalf("%%p leaks the secret: %q", out)
	}
}

func TestMasksWhenEmbeddedInStruct(t *testing.T) {
	type credentials struct {
		User     string
		Password Secret
	}
	cfg := credentials{User: "app", Password: New(value)}

	if out := fmt.Sprintf("%+v", cfg); strings.Contains(out, value) {
		t.Fatalf("formatted struct leaks the secret: %s", out)
	}
}

func TestMasksStandardSerialization(t *testing.T) {
	type credentials struct {
		Password Secret `json:"password" xml:"password"`
	}

	s := New(value)

	jsonData, err := json.Marshal(credentials{Password: s})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	assertMaskedBytes(t, "json", jsonData)

	xmlData, err := xml.Marshal(credentials{Password: s})
	if err != nil {
		t.Fatalf("xml.Marshal: %v", err)
	}
	assertMaskedBytes(t, "xml", xmlData)

	textData, err := s.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	assertMaskedBytes(t, "text", textData)

	var gobData bytes.Buffer
	if err := gob.NewEncoder(&gobData).Encode(s); err == nil {
		assertDoesNotLeak(t, "gob", gobData.Bytes())
	}
}

func TestMasksInSlog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	logger.Info("connecting", "password", New(value))

	assertMaskedBytes(t, "slog", buf.Bytes())
}

func TestUnmarshalTextReadsInButOutputsStayMasked(t *testing.T) {
	var s Secret
	if err := s.UnmarshalText([]byte(value)); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}
	if got := s.Reveal(); got != value {
		t.Fatalf("Reveal() = %q, want original value", got)
	}
	if out := fmt.Sprintf("%+v", s); strings.Contains(out, value) {
		t.Fatalf("formatting leaks value after UnmarshalText: %q", out)
	}
}

func TestEmptyOutputsStayEmpty(t *testing.T) {
	var s Secret
	if got := s.String(); got != "" {
		t.Errorf("empty String() = %q", got)
	}
	if got := fmt.Sprintf("%v", s); got != "" {
		t.Errorf("empty fmt output = %q", got)
	}
	data, err := s.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("empty MarshalText() = %q", data)
	}
}

func assertMaskedBytes(t *testing.T, name string, data []byte) {
	t.Helper()
	assertDoesNotLeak(t, name, data)
	if !bytes.Contains(data, []byte(mask)) {
		t.Errorf("%s output is not masked: %q", name, data)
	}
}

func assertDoesNotLeak(t *testing.T, name string, data []byte) {
	t.Helper()
	if bytes.Contains(data, []byte(value)) {
		t.Fatalf("%s leaks the secret: %q", name, data)
	}
}
