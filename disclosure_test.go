package secret_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/uchaloop/secret/v2"
)

func TestNestedValuesDoNotDiscloseSecret(t *testing.T) {
	const raw = "private-test-password"
	password := secret.New(raw)
	values := []any{
		struct{ Password secret.Secret }{password},
		struct{ password secret.Secret }{password},
		struct{ password *secret.Secret }{&password},
		struct{ credentials []secret.Secret }{[]secret.Secret{password}},
		struct{ credentials map[string]secret.Secret }{map[string]secret.Secret{"password": password}},
		struct{ credentials any }{password},
		struct {
			private struct{ Password secret.Secret }
		}{struct{ Password secret.Secret }{password}},
	}

	for _, value := range values {
		for _, format := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%X", "%d", "%f", "%t", "%p", "%w", "%20.3s", "%#q"} {
			output := fmt.Sprintf(format, value)
			if strings.Contains(output, raw) || strings.Contains(strings.ToLower(output), fmt.Sprintf("%x", raw)) {
				t.Errorf("%T with %s disclosed the secret", value, format)
			}
		}

		var output bytes.Buffer

		for _, handler := range []slog.Handler{
			slog.NewTextHandler(&output, nil),
			slog.NewJSONHandler(&output, nil),
		} {
			output.Reset()
			slog.New(handler).Info("config", "credentials", value)

			if strings.Contains(output.String(), raw) {
				t.Errorf("%T disclosed a secret in %T", handler, value)
			}
		}
	}
}

func TestCopiesKeepTheirValues(t *testing.T) {
	original := secret.New("original")
	snapshot := original
	original.Clear()

	if !original.IsZero() || snapshot.Reveal() != "original" {
		t.Fatal("clearing a secret changed its copy")
	}

	original = snapshot
	if err := original.UnmarshalText([]byte("replacement")); err != nil {
		t.Fatal(err)
	}

	if original.Reveal() != "replacement" || snapshot.Reveal() != "original" {
		t.Fatal("replacing a secret changed its copy")
	}
}

func TestSecretIsNotComparable(t *testing.T) {
	if reflect.TypeOf(secret.Secret{}).Comparable() {
		t.Fatal("Secret must not support equality or map keys")
	}
}

func TestEmptyInputRestoresZeroValue(t *testing.T) {
	empty := secret.New("")
	if !empty.IsZero() || !reflect.ValueOf(empty).IsZero() || empty.Reveal() != "" {
		t.Fatal("empty construction must preserve the zero value")
	}

	password := secret.New("initial")
	if err := password.UnmarshalText(nil); err != nil {
		t.Fatal(err)
	}

	if !password.IsZero() || !reflect.ValueOf(password).IsZero() || password.Reveal() != "" {
		t.Fatal("empty input must restore the zero value")
	}

	password = secret.New("initial")
	password.Clear()

	if !reflect.ValueOf(password).IsZero() {
		t.Fatal("Clear must restore the zero value")
	}
}
