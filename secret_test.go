package secret

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

const value = "hunter2"

func TestReveal(t *testing.T) {
	if got := Secret(value).Reveal(); got != value {
		t.Fatalf("Reveal() = %q, want %q", got, value)
	}
}

func TestMasksAllTextualForms(t *testing.T) {
	s := Secret(value)

	forms := map[string]string{
		"%v":  fmt.Sprintf("%v", s),
		"%s":  fmt.Sprintf("%s", s),
		"%+v": fmt.Sprintf("%+v", s),
		"%#v": fmt.Sprintf("%#v", s),
	}
	for verb, out := range forms {
		if strings.Contains(out, value) {
			t.Errorf("%s leaks the secret: %q", verb, out)
		}
		if !strings.Contains(out, mask) {
			t.Errorf("%s = %q, want it to contain the mask", verb, out)
		}
	}
}

func TestMasksInStruct(t *testing.T) {
	type creds struct {
		User     string
		Password Secret
	}
	c := creds{User: "app", Password: Secret(value)}

	if out := fmt.Sprintf("%+v", c); strings.Contains(out, value) {
		t.Fatalf("struct %%+v leaks the secret: %s", out)
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if bytes.Contains(data, []byte(value)) {
		t.Fatalf("json leaks the secret: %s", data)
	}
	if !bytes.Contains(data, []byte(mask)) {
		t.Fatalf("json = %s, want it to contain the mask", data)
	}
}

func TestMasksInSlog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	logger.Info("connecting", "password", Secret(value))

	if strings.Contains(buf.String(), value) {
		t.Fatalf("slog leaks the secret: %s", buf.String())
	}
	if !strings.Contains(buf.String(), mask) {
		t.Fatalf("slog output = %q, want it to contain the mask", buf.String())
	}
}

func TestUnmarshalTextReadsInButOutputStaysMasked(t *testing.T) {
	var s Secret
	if err := s.UnmarshalText([]byte(value)); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}

	if got := s.Reveal(); got != value {
		t.Fatalf("Reveal() = %q, want %q", got, value)
	}
	if out := s.String(); strings.Contains(out, value) {
		t.Fatalf("String() leaks after UnmarshalText: %q", out)
	}
}

func TestEmptyStaysEmpty(t *testing.T) {
	var s Secret

	if got := s.String(); got != "" {
		t.Errorf("empty String() = %q, want empty", got)
	}
	data, err := s.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("empty MarshalText() = %q, want empty", data)
	}
}
