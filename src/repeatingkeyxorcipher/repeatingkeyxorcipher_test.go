package repeatingkeyxorcipher

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// 5. Repeating-key XOR Cipher
func TestMetasanoRepeatingKeyXORCipher(t *testing.T) {
	const input = "Burning 'em, if you ain't quick and nimble\n" +
		"I go crazy when I hear a cymbal"
	const key = "ICE"
	const expected = "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"

	inbytes, _ := hex.DecodeString(fmt.Sprintf("%x", input))
	keybytes, _ := hex.DecodeString(fmt.Sprintf("%x", key))

	out := EncodeRepeatingKeyXORCipher(inbytes, keybytes)
	if hex.EncodeToString(out) != expected {
		t.Errorf("%v (len %v)", out, len(out))
	}
}
