/*
3. Single-character XOR Cipher

The hex encoded string:

   1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736

... has been XOR'd against a single character. Find the key, decrypt
the message.

Write code to do this for you. How? Devise some method for "scoring" a
piece of English plaintext. (Character frequency is a good metric.)
Evaluate each output and choose the one with the best score.

Tune your algorithm until this works.
*/

package singlecharxorcipher

import "errors"

// DecryptSingleCharXORCipher finds the best
// English plaintext alternative and returns that
// together with the score of the alternative or
// and error. The closer the score is to the
// ciphertext length the better.
// TODO: handle ties and use the unicode library
//       RangeTables and In() function for scoring
func DecryptSingleCharXORCipher(ciphertext []byte) (string, int, error) {
	cipherlen := len(ciphertext)
	topCandidate := -1
	topScore := -1
	for i := 0; i < 256; i++ {
		testbyte := byte(i)
		score := 0
		for j := 0; j < cipherlen; j++ {
			out := byte(ciphertext[j] ^ testbyte)
			if byte(0x30) <= out && out <= byte(0x39) { // 0-9
				score += 1
			}
			// ...English plain-text, i.e. forget unicode printable for now
			if byte(0x61) <= out && out <= byte(0x7a) { // a-z
				score += 1
			}
			if byte(0x41) <= out && out <= byte(0x5a) { // A-Z
				score += 1
			}
			if out == byte(0x2e) || out == byte(0x2c) || out == byte(0x20) { // . , SP
				score += 1
			}
		}
		if topScore < score {
			topScore = score
			topCandidate = i
		}
	}
	if topCandidate < 0 {
		return "", -1, errors.New("could not detect a proper candidate")
	}
	result := make([]byte, cipherlen)
	topCandidateByte := byte(topCandidate)
	for i := 0; i < cipherlen; i++ {
		result[i] = byte(ciphertext[i] ^ topCandidateByte)
	}
	return string(result), topScore, nil
}
