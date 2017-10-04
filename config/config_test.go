package config

import (
	"testing"
)

func TestKeyLength(t *testing.T) {
	validKey := "00000000000000000000000000000000"
	byteKey, err := validateKeyLength(validKey)
	if err != nil {
		t.Errorf("unable to validate simple key (byteKey = [%s]):\n%s\n",
			err, byteKey)
	}
	if len(validKey) != len(byteKey) {
		t.Errorf("keys were reported valid, but are not the same size")
	}

	invalidKey := "000000"
	byteKey, err = validateKeyLength(invalidKey)
	if err == nil {
		t.Errorf("no error thrown on invalid key")
	}

	t.Logf("key length validated")
}
