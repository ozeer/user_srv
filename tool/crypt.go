package tool

import (
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

const KEY = "dfertf12dfertf12"

// 加密
func Encrypt(key string, plainText string) string {
	crypt := crypto.
		FromString(plainText).
		SetKey(key).
		Aes().
		ECB().
		PKCS7Padding().
		Encrypt().
		ToBase64String()

	return crypt
}

// 解密
func Decrypt(key string, encryptText string) string {
	cyptde := crypto.
		FromBase64String(encryptText).
		SetKey(key).
		Aes().
		ECB().
		PKCS7Padding().
		Decrypt().
		ToString()

	return cyptde
}
