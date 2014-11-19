package aesecb

import (
	"../base64"
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

/*
7. AES in ECB Mode

The Base64-encoded content at the following location:

 https://gist.github.com/3132853

Has been encrypted via AES-128 in ECB mode under the key

 "YELLOW SUBMARINE".

(I like "YELLOW SUBMARINE" because it's exactly 16 bytes long).

Decrypt it.

Easiest way:

Use OpenSSL::Cipher and give it AES-128-ECB as the cipher.
*/

func TestDecryptAES128ECB(t *testing.T) {
	const keyString = "YELLOW SUBMARINE"
	const expected = "Play that funky music white boy"

	file, err := os.Open("gistfile1.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close() // todo: closing can fail

	lineScanner := bufio.NewScanner(file)

	base64CipherText := ""
	for lineScanner.Scan() {
		base64CipherText += lineScanner.Text()
	}

	base64CipherBytes := []byte(base64CipherText)
	cipherBytes := base64.ConvertFromBase64(base64CipherBytes)

	plainBytes := DecryptAES128ECB(cipherBytes, []byte(keyString))

	if len(cipherBytes) != len(plainBytes) {
		t.Fatal("length mismatch")
	}
	if !strings.Contains(fmt.Sprintf("%s", plainBytes), expected) {
		t.Fatal("could not find expected plaintext string")
	}
}
