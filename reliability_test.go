package secret_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/uchaloop/secret/v2"
)

func TestTextInputOwnsItsBytes(t *testing.T) {
	input := []byte("original")
	var value secret.Secret
	if err := value.UnmarshalText(input); err != nil {
		t.Fatal(err)
	}

	input[0] = 'X'
	if value.Reveal() != "original" {
		t.Fatal("modifying input bytes changed the secret")
	}
}

func TestConcurrentReadsAndIndependentCopies(t *testing.T) {
	shared := secret.New("private-input")
	var workers sync.WaitGroup

	for workerIndex := 0; workerIndex < 16; workerIndex++ {
		workers.Add(1)

		go func() {
			defer workers.Done()

			for iteration := 0; iteration < 100; iteration++ {
				if shared.Reveal() != "private-input" || shared.IsZero() || fmt.Sprint(shared) != "****" {
					t.Error("concurrent reading changed the secret")

					return
				}

				data, err := json.Marshal(shared)
				if err != nil {
					t.Error(err)

					return
				}

				if string(data) != `"****"` {
					t.Error("concurrent serialization failed")

					return
				}

				local := shared
				local.Clear()
				if err := local.UnmarshalText([]byte("replacement")); err != nil {
					t.Error(err)

					return
				}
			}
		}()
	}

	workers.Wait()

	if shared.Reveal() != "private-input" {
		t.Fatal("changing a copy changed the shared secret")
	}
}

func FuzzSecretTextMasking(f *testing.F) {
	for _, seed := range []string{"", "password", "\x00\xff", "\"\\\n", "пароль", "****"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, raw string) {
		var value secret.Secret
		if err := value.UnmarshalText([]byte(raw)); err != nil {
			t.Fatal(err)
		}

		if value.Reveal() != raw {
			t.Fatal("text input changed")
		}

		expected := "****"
		if len(raw) == 0 {
			expected = ""
		}

		if fmt.Sprintf("%v", value) != expected || fmt.Sprintf("%q", value) != fmt.Sprintf("%q", expected) {
			t.Fatal("unexpected formatted representation")
		}

		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}

		if string(encoded) != fmt.Sprintf("%q", expected) {
			t.Fatal("unexpected JSON representation")
		}

		snapshot := value
		value.Clear()

		if !value.IsZero() || snapshot.Reveal() != raw {
			t.Fatal("Clear changed a copy")
		}
	})
}
