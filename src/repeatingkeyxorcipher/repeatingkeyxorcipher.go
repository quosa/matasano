/*
5. Repeating-key XOR Cipher

Write the code to encrypt the string:

Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal

Under the key "ICE", using repeating-key XOR. It should come out to:

0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f

Encrypt a bunch of stuff using your repeating-key XOR function. Get a
feel for it.
*/

package repeatingkeyxorcipher

import "errors"

func EncodeRepeatingKeyXORCipher(inbytes, keybytes []byte) []byte {
	inlen := len(inbytes)
	keylen := len(keybytes)
	result := make([]byte, inlen)
	for i := 0; i < inlen; i++ {
		result[i] = byte(inbytes[i] ^ keybytes[i%keylen])
	}
	return result
}

/*
6. Break repeating-key XOR

The buffer at the following location:

https://gist.github.com/3132752

is base64-encoded repeating-key XOR. Break it.

Here's how:

a. Let KEYSIZE be the guessed length of the key; try values from 2 to
(say) 40.

b. Write a function to compute the edit distance/Hamming distance
between two strings. The Hamming distance is just the number of
differing bits. The distance between:

this is a test

and:

wokka wokka!!!

is 37.

c. For each KEYSIZE, take the FIRST KEYSIZE worth of bytes, and the
SECOND KEYSIZE worth of bytes, and find the edit distance between
them. Normalize this result by dividing by KEYSIZE.

d. The KEYSIZE with the smallest normalized edit distance is probably
the key. You could proceed perhaps with the smallest 2-3 KEYSIZE
values. Or take 4 KEYSIZE blocks instead of 2 and average the
distances.

e. Now that you probably know the KEYSIZE: break the ciphertext into
blocks of KEYSIZE length.

f. Now transpose the blocks: make a block that is the first byte of
every block, and a block that is the second byte of every block, and
so on.

g. Solve each block as if it was single-character XOR. You already
have code to do this.

e. For each block, the single-byte XOR key that produces the best
looking histogram is the repeating-key XOR key byte for that
block. Put them together and you have the key.
*/

// Calculate the Hamming distance between 2 byte arrays.
// Returns -1 distance in error cases.
// http://en.wikipedia.org/wiki/Hamming_distance for info
func CalculateDistance(input1, input2 []byte) (int, error) {
	len1 := len(input1)
	if len1 != len(input2) {
		return -1, errors.New("buffers must have equal length!")
	}

	result := make([]byte, len1)
	distance := 0

	for i := 0; i < len1; i++ {
		result[i] = input1[i] ^ input2[i]
		masks := []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
		for _, m := range masks { // TODO: use left-shifts
			if result[i]&byte(m) > 0 {
				distance += 1
			}
		}
	}
	return distance, nil
}
