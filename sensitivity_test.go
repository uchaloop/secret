package secret_test

import (
	"fmt"
	"github.com/uchaloop/secret/v2"
	"testing"
)

func TestSensitiveMarkerCompatibility(t *testing.T) {
	s := secret.New("do-not-disclose")
	marker, ok := any(s).(interface{ IsSensitive() })
	if !ok {
		t.Fatal("Secret must expose the sensitive marker")
	}
	marker.IsSensitive()
	var _ secret.Value = s
	if s.Reveal() != "do-not-disclose" || fmt.Sprint(s) != "****" {
		t.Fatal("marker changed secret behavior")
	}
}
