package encrypt

import (
	"fmt"
	"testing"
)

var (
	key = []byte("1234567890123456")
)

func TestEncrypt(t *testing.T) {
	var err error
	src := "abcabc"
	cy, err := AesEcryptCBC([]byte(src), key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cy)
}

func TestAesDeCryptCBC(t *testing.T) {
	cy := "3iHrBvI1IW2X7WovlWKWmg=="
	src, err := AesDeCryptCBC(cy, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(src)
}
