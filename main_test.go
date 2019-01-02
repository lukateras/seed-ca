package main

import (
	"encoding/base64"
	"golang.org/x/crypto/blake2b"
	"io/ioutil"
	"os"
	"testing"
)

func testFile(t *testing.T, name, expectedHash string) {
	f, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	hashBytes := blake2b.Sum512(f)
	hash := base64.RawStdEncoding.EncodeToString(hashBytes[:])

	if hash != expectedHash {
		t.Errorf(`
%s
  expected: %s
  got: %s`, name, expectedHash, hash)
	}
}

func Test(t *testing.T) {
	err := os.Setenv("SEED_CA", "")
	if err != nil {
		t.Fatal(err)
	}

	main()

	testFile(t, "ca.key", "ZDWPrBOVoQBv4u6vruq0rXo72vPnH8pJNIVwRwQQcwZHpWbucz0H5QD5JYP+TtZiXF36pt1/bp9LLLpCktZrdw")
	testFile(t, "ca.crt", "Lhm6DUqgSALS3LxMSQ+AgQDBaNmDImkbArM91aeWjhsgUTGe6cBJ+f8Jzj8TwlQbtNEts6tyQ2E0tgmoE6DTZA")
}
