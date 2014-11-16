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

func EncodeRepeatingKeyXORCipher(inbytes, keybytes []byte) []byte {
	inlen := len(inbytes)
	keylen := len(keybytes)
	result := make([]byte, inlen)
	for i := 0; i < inlen; i++ {
		result[i] = byte(inbytes[i] ^ keybytes[i%keylen])
	}
	return result
}