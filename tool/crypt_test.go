package tool

import (
	"fmt"
	"testing"
)

func TestCrypt(t *testing.T) {
	secret := Encrypt(KEY, "hello")
	fmt.Println(secret)
}

func TestDeCrypt(t *testing.T) {
	plainText := Decrypt(KEY, "ZVB50oXjc2vW3n20er/nbg==")
	fmt.Println(plainText)
}
