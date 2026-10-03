package secret_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"testing"

	"github.com/uchaloop/secret/v2"
)

var (
	_ secret.Value = secret.Secret{}
	_ secret.Value = (*secret.Secret)(nil)
)

func TestNewRevealAndZero(t *testing.T) {
	var zero secret.Secret
	if !zero.IsZero() || zero.Reveal() != "" {
		t.Fatal("zero Secret must be empty")
	}

	value := secret.New("password")
	if value.IsZero() || value.Reveal() != "password" {
		t.Fatal("New must preserve the supplied value")
	}
}

func TestClearNilReceiver(t *testing.T) {
	var value *secret.Secret
	value.Clear()
}

func TestMaskedOutput(t *testing.T) {
	for _, raw := range []string{"", "private-password"} {
		t.Run(raw, func(t *testing.T) {
			value := secret.New(raw)
			expected := "****"
			if len(raw) == 0 {
				expected = ""
			}

			if value.String() != expected || value.GoString() != fmt.Sprintf("%q", expected) {
				t.Fatal("unexpected string representation")
			}

			for _, format := range []string{"%v", "%s", "%+v", "%#v", "%x", "%X", "%d"} {
				if fmt.Sprintf(format, value) != expected {
					t.Errorf("unexpected output for %s", format)
				}
			}

			text, err := value.MarshalText()
			if err != nil {
				t.Fatal(err)
			}

			if string(text) != expected {
				t.Fatal("unexpected text representation")
			}

			jsonData, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}

			if string(jsonData) != fmt.Sprintf("%q", expected) {
				t.Fatal("unexpected JSON representation")
			}

			xmlData, err := xml.Marshal(struct {
				XMLName  xml.Name      `xml:"credentials"`
				Password secret.Secret `xml:"password"`
			}{Password: value})
			if err != nil {
				t.Fatal(err)
			}

			if string(xmlData) != "<credentials><password>"+expected+"</password></credentials>" {
				t.Fatal("unexpected XML representation")
			}

			if value.LogValue().String() != expected {
				t.Fatal("unexpected slog representation")
			}

			var output bytes.Buffer

			for _, handler := range []slog.Handler{
				slog.NewTextHandler(&output, nil),
				slog.NewJSONHandler(&output, nil),
			} {
				output.Reset()
				slog.New(handler).Info("config", "password", value)

				if len(raw) > 0 && (bytes.Contains(output.Bytes(), []byte(raw)) || !bytes.Contains(output.Bytes(), []byte(expected))) {
					t.Fatalf("%T did not mask the secret", handler)
				}
			}
		})
	}
}
