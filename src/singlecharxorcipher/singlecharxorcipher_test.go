package singlecharxorcipher

import (
	"bufio"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

// 3. Single-character XOR Cipher
func TestMetasanoDecryptSingleCharXORCipher(t *testing.T) {
	const ciphertext = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

	cipherbytes, _ := hex.DecodeString(ciphertext)

	out, score, err := DecryptSingleCharXORCipher(cipherbytes)
	if err != nil || score < 0 {
		t.Fatal(err)
	}
	if out != "Cooking MC's like a pound of bacon" {
		t.Errorf("%v (len %v)", out, len(out))
	}
}

/*
4. Detect single-character XOR

One of the 60-character strings at:

https://gist.github.com/3132713

has been encrypted by single-character XOR. Find it. (Your code from
#3 should help.)
*/
func TestMetasanoFindSingleCharXORCipher(t *testing.T) {
	file, err := os.Open("gistfile1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close() // todo: closing can fail

	lineScanner := bufio.NewScanner(file)

	topScore := -1
	topResult := ""
	for lineScanner.Scan() {
		cipherbytes, _ := hex.DecodeString(lineScanner.Text())
		out, score, err := DecryptSingleCharXORCipher(cipherbytes)
		if err != nil || score < 0 {
			t.Fatal(out, err)
		}
		if score > topScore {
			topScore = score
			topResult = out
		}
	}
	if topScore < 0 || topResult == "" {
		t.Fatal("could not find plaintext")
	}
	if strings.TrimSpace(topResult) != "Now that the party is jumping" {
		t.Errorf("%v (len %v)", topResult, len(topResult))
	}
}
