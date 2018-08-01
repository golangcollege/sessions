package sessions

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key1 := [32]byte{}
	copy(key1[:], []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"))

	key2 := [32]byte{}
	copy(key2[:], []byte("3j4a0lniSrNb4xMdkYjsgG74mjRCF75u"))

	message1 := []byte("foo bar baz")
	token, err := encrypt(message1, key2)
	if err != nil {
		t.Fatal(err)
	}

	message2, err := decrypt(token, [][32]byte{key1, key2})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(message1, message2) {
		t.Errorf("got %q: expect %q", message2, message1)
	}

	_, err = decrypt(token, [][32]byte{key1})
	if err != errInvalidToken {
		t.Errorf("got %v: expect %q", err, errInvalidToken)
	}
}
