package aesecb

import "crypto/aes"

// Go left out AES ECB because it is insecure
// https://code.google.com/p/go/issues/detail?id=5597
func DecryptAES128ECB(cipher, key []byte) []byte {
	if len(cipher)%16 != 0 || len(key) != 16 {
		panic("wrong length")
	}

	cb, _ := aes.NewCipher(key)
	plain := make([]byte, 0)
	block := make([]byte, 16)
	for len(cipher) > 0 {
		cb.Decrypt(block, cipher[:16])
		cipher = cipher[16:]
		plain = append(plain, block...)
	}
	return plain
}
