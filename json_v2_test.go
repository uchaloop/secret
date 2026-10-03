//go:build go1.27

package secret_test

import (
	json "encoding/json/v2"
	"testing"

	"github.com/uchaloop/secret/v2"
)

func TestJSONV2Input(t *testing.T) {
	checkJSONDecodeContract(t, func(data []byte, value *secret.Secret) error {
		return json.Unmarshal(data, value)
	})
}

func TestJSONV2OutputStaysMasked(t *testing.T) {
	for _, test := range []struct {
		value secret.Secret
		want  string
	}{
		{secret.New("private-input"), `"****"`},
		{secret.Secret{}, `""`},
	} {
		data, err := json.Marshal(test.value)
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != test.want {
			t.Fatal("unexpected JSON representation")
		}
	}
}
