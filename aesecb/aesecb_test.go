package aesecb

import (
	"../base64"
	"bufio"
	"encoding/hex"
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

/*
8. Detecting ECB

At the following URL are a bunch of hex-encoded ciphertexts:

https://gist.github.com/3132928

One of them is ECB encrypted. Detect it.

Remember that the problem with ECB is that it is stateless and
deterministic; the same 16 byte plaintext block will always produce
the same 16 byte ciphertext.
*/
func TestDetectAES128ECB(t *testing.T) {
	const expected = "d880619740a8a19b7840a8a31c810a3d08649af70dc06f4fd5d2d69c744cd283e2dd052f6b641dbf9d11b0348542bb5708649af70dc06f4fd5d2d69c744cd2839475c9dfdbc1d46597949d9c7e82bf5a08649af70dc06f4fd5d2d69c744cd28397a93eab8d6aecd566489154789a6b0308649af70dc06f4fd5d2d69c744cd283d403180c98c8f6db1f2a3f9c4040deb0ab51b29933f2c123c58386b06fba186a"
	file, err := os.Open("gistfile2.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close() // todo: closing can fail

	lineScanner := bufio.NewScanner(file)

	found := ""
	for lineScanner.Scan() {
		cipherHex := lineScanner.Text()
		cipherBytes, _ := hex.DecodeString(cipherHex)
		if DetectAES128ECBCipher(cipherBytes) {
			found = cipherHex
			fmt.Printf("This looks like AES 128 ECB cipher: %s\n", cipherHex)
		}
	}
	if found == "" {
		t.Fatal("Couldn't find an AES 128 ECB cipher string")
	}
	if found != expected {
		t.Fatal("Detected the wrong cipher")
	}
}
