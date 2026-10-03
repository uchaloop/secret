package secret_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/uchaloop/secret/v2"
)

func checkJSONDecodeContract(t *testing.T, decode func([]byte, *secret.Secret) error) {
	t.Helper()

	for _, test := range []struct{ input, want string }{
		{`"replacement"`, "replacement"},
		{`""`, ""},
		{`null`, ""},
		{" \n null\t", ""},
		{`"a\n\"b\u263a"`, "a\n\"b☺"},
		{`"\ud83d\ude00"`, "😀"},
		{`"\uD800\uDC00\uDBFF\uDFFF"`, "𐀀\U0010ffff"},
		{`"\ufffd�"`, "��"},
		{`"\\ud800"`, `\ud800`},
	} {
		t.Run(test.input, func(t *testing.T) {
			value := secret.New("original")
			snapshot := value
			if err := decode([]byte(test.input), &value); err != nil {
				t.Fatal(err)
			}

			if value.Reveal() != test.want || value.IsZero() != (len(test.want) == 0) {
				t.Fatal("unexpected decoded value")
			}

			if snapshot.Reveal() != "original" {
				t.Fatal("decoding changed an existing snapshot")
			}
		})
	}

	for _, input := range []string{
		`true`,
		`42`,
		`[]`,
		`{"password":"private-input"}`,
		`"unterminated`,
		``,
		"\u00a0null\u00a0",
		"\vnull\f",
		"\"\xff\"",
		`"\ud800"`,
		`"\udc00"`,
		`"\ud800x"`,
		`"\ud800\ud800"`,
		`"\ud800\u0041"`,
		`"\ud800\\udc00"`,
	} {
		value := secret.New("original")
		if err := decode([]byte(input), &value); err == nil {
			t.Fatal("invalid JSON input was accepted")
		}

		if value.Reveal() != "original" {
			t.Fatal("failed decoding changed the secret")
		}
	}
}

func TestJSONInput(t *testing.T) {
	checkJSONDecodeContract(t, func(data []byte, value *secret.Secret) error {
		return json.Unmarshal(data, value)
	})
}

func TestDirectJSONInput(t *testing.T) {
	checkJSONDecodeContract(t, func(data []byte, value *secret.Secret) error {
		err := value.UnmarshalJSON(data)
		if err == nil {
			return nil
		}

		if err.Error() != "secret: expected a JSON string or null" || errors.Unwrap(err) != nil {
			t.Fatal("expected a safe error without input or a decoder cause")
		}

		return err
	})

	for _, input := range []string{`"private-input`, `private-input`, `null true`, `"one" "two"`} {
		value := secret.New("original")
		err := value.UnmarshalJSON([]byte(input))
		if err == nil || err.Error() != "secret: expected a JSON string or null" || errors.Unwrap(err) != nil {
			t.Fatal("expected a safe error without input or a decoder cause")
		}

		if value.Reveal() != "original" {
			t.Fatal("failed direct decoding changed the secret")
		}
	}
}

func TestNilDecodeReceiver(t *testing.T) {
	var value *secret.Secret
	if err := value.UnmarshalJSON([]byte(`null`)); err == nil {
		t.Fatal("nil JSON receiver must fail")
	}

	if err := value.UnmarshalText([]byte("input")); err == nil {
		t.Fatal("nil text receiver must fail")
	}
}

func FuzzJSONInput(f *testing.F) {
	for _, input := range []string{`null`, `"password"`, `"\ud800"`, `"\ud83d\ude00"`, `"\\ud800"`, "\u00a0null", "\"\xff\""} {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, input string) {
		value := secret.New("original")
		err := value.UnmarshalJSON([]byte(input))
		if err != nil {
			if value.Reveal() != "original" || err.Error() != "secret: expected a JSON string or null" || errors.Unwrap(err) != nil {
				t.Fatal("failed decoding must preserve the secret and return a safe error")
			}

			return
		}

		var expected *string
		if err := json.Unmarshal([]byte(input), &expected); err != nil {
			t.Fatal("accepted invalid JSON")
		}

		if expected == nil {
			if !value.IsZero() {
				t.Fatal("null must clear the secret")
			}

			return
		}

		if value.Reveal() != *expected {
			t.Fatal("decoded string differs from input")
		}
	})
}
