package encrypt

import (
	"testing"
)

var aesKey = "aDgeDGhe5gh512hgaDgeDGhe5gh512hg"

func TestAesEncode(t *testing.T) {
	var (
		corpID    = "accc"
		plaintext = "abcd一哈哈汪畅叙1add在-=+%6$#@&*!"
	)
	ciphertext, err := NewAES([]byte(aesKey), WithAesSalt(corpID)).EncryptString(plaintext)
	if err != nil {
		t.Error(err)
	}
	t.Log(ciphertext)
}

func TestAesDecode(t *testing.T) {
	var (
		corpID     = "accc"
		ciphertext = "IeDUTgDuMGuIRLmie7YHigLramOqDDbLiDW2SqLUDa0qpEeyB9ghFrWNxnGcAwiiT3mXT/8o1YoGyKVM58tkEfiAJuUCgIu9ax9UuHnFWKQ="
	)
	plaintext, err := NewAES([]byte(aesKey), WithAesSalt(corpID)).DecryptString(ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(plaintext)
}

func TestRandomString(t *testing.T) {
	t.Log(randomString(16))
}
