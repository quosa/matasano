package repeatingkeyxorcipher

import (
	"../base64"
	"../singlecharxorcipher"
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strings"
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

// This test is for debugging assignment
// 6. "Break repeating-key XOR".
// Here we create the flow manually so that we know
// where the process goes wrong.
// Running the transpose guess part manually,
// the fault is in the keysize guess part
// and return too high normalized distance
// for the correct key size.
func TestRepeatingKeyXORCipherForBreaking(t *testing.T) {
	const plain = "Burning 'em, if you ain't quick and nimble\n" +
		"I go crazy when I hear a cymbal"
	const key = "ICE"
	const expected = "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"

	inbytes, _ := hex.DecodeString(fmt.Sprintf("%x", plain))
	keybytes, _ := hex.DecodeString(fmt.Sprintf("%x", key))

	out := EncodeRepeatingKeyXORCipher(inbytes, keybytes)
	if hex.EncodeToString(out) != expected {
		t.Errorf("%v (len %v)", out, len(out))
	}

	base64bytes := base64.ConvertToBase64(out)

	// Starting to backtrack
	cipherBytes := base64.ConvertFromBase64(base64bytes)

	// test the keysize guessing
	leastDistance := float32(math.MaxFloat32)
	leastKeysize := -1
	for keysize := 2; keysize < 20; keysize++ {
		sliceA := cipherBytes[:keysize]
		sliceB := cipherBytes[keysize : 2*keysize]
		sliceC := cipherBytes[2*keysize : 3*keysize]
		sliceD := cipherBytes[3*keysize : 4*keysize]
		distance1, _ := CalculateDistance(sliceA, sliceB)
		distance2, _ := CalculateDistance(sliceB, sliceC)
		distance3, _ := CalculateDistance(sliceC, sliceD)

		normalizedDistance := float32(distance1+distance2+distance3) / float32(3*keysize)
		fmt.Printf("keysize %2v gives dist %v\n", keysize, normalizedDistance)

		if normalizedDistance < leastDistance {
			leastDistance = normalizedDistance
			leastKeysize = keysize
		}
	}
	fmt.Printf("found keysize %v with dist %v\n", leastKeysize, leastDistance)

	// --> gives high number for the correct value
	// NOTE: lock the correct keysize to test the remaining part
	//       and fail the test only at the end
	if leastKeysize != 3 {
		t.Error("failed to find the correct keysize 3, got", leastKeysize)
		leastKeysize = 3
	}

	transpose := make(map[int][]byte)
	for i, v := range cipherBytes {
		transpose[i%leastKeysize] = append(transpose[i%leastKeysize], v)
	}

	repeatingKey := make([]byte, leastKeysize)
	for i := 0; i < leastKeysize; i++ {
		key, score, err := singlecharxorcipher.BreakSingleCharXORCipher(transpose[i])
		if err != nil || score < 0 {
			t.Fatal(score, err)
		}
		repeatingKey[i] = key
	}

	if !bytes.Equal(repeatingKey, keybytes) {
		t.Fatal("key mismatch")
	}

	plainBytes := EncodeRepeatingKeyXORCipher(cipherBytes, repeatingKey)
	plainText := fmt.Sprintf("%s", plainBytes)
	if plain != plainText {
		t.Fatal("plain text mismatch")
	}
}

// 6. Break repeating-key XOR
func TestMetasanoBreakRepeatingKeyXORCipher(t *testing.T) {
	const expected = "Play that funky music white boy"

	file, err := os.Open("gistfile1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close() // todo: closing can fail

	lineScanner := bufio.NewScanner(file)

	// is base64-encoded repeating-key XOR
	base64CipherText := ""
	for lineScanner.Scan() {
		base64CipherText += lineScanner.Text()
	}

	base64CipherBytes, err := hex.DecodeString(fmt.Sprintf("%x", base64CipherText))
	if err != nil {
		t.Fatal(err)
	}
	cipherBytes := base64.ConvertFromBase64(base64CipherBytes)

	leastDistance := float32(math.MaxFloat32)
	leastKeysize := -1

	for keysize := 2; keysize < 40; keysize++ {
		sliceA := cipherBytes[:keysize]
		sliceB := cipherBytes[keysize : 2*keysize]
		sliceC := cipherBytes[2*keysize : 3*keysize]
		sliceD := cipherBytes[3*keysize : 4*keysize]
		distance1, _ := CalculateDistance(sliceA, sliceB)
		distance2, _ := CalculateDistance(sliceB, sliceC)
		distance3, _ := CalculateDistance(sliceC, sliceD)

		normalizedDistance := float32(distance1+distance2+distance3) / float32(3*keysize)
		fmt.Printf("keysize %v gives dist %v\n", keysize, normalizedDistance)

		if normalizedDistance < leastDistance {
			leastDistance = normalizedDistance
			leastKeysize = keysize
		}
	}
	if leastKeysize < 0 {
		t.Fatal("could not find keysize candidate")
	}
	fmt.Printf("found keysize %v with dist %v\n", leastKeysize, leastDistance)

	// NOTE: lock the correct keysize to test the remaining part
	//       and fail the test only at the end
	if leastKeysize != 29 {
		t.Error("failed to find the correct keysize 29, got", leastKeysize)
		leastKeysize = 29
	}

	//e. Now that you probably know the KEYSIZE: break the ciphertext into
	//blocks of KEYSIZE length.
	//f. Now transpose the blocks: make a block that is the first byte of
	//every block, and a block that is the second byte of every block, and
	//so on.
	transpose := make(map[int][]byte)
	for i, v := range cipherBytes {
		transpose[i%leastKeysize] = append(transpose[i%leastKeysize], v)
	}

	//g. Solve each block as if it was single-character XOR.
	//You already have code to do this.
	repeatingKey := make([]byte, leastKeysize)
	for i := 0; i < leastKeysize; i++ {
		key, score, err := singlecharxorcipher.BreakSingleCharXORCipher(transpose[i])
		if err != nil || score < 0 {
			t.Fatal(score, err)
		}
		repeatingKey[i] = key
	}

	//e. For each block, the single-byte XOR key that produces the best
	//looking histogram is the repeating-key XOR key byte for that
	//block. Put them together and you have the key.
	plainBytes := EncodeRepeatingKeyXORCipher(cipherBytes, repeatingKey)

	if !strings.Contains(fmt.Sprintf("%s", plainBytes), expected) {
		t.Fatal("could not find expected plaintext string")
	}

}

// 6. Break repeating-key XOR
// Tests for the Hamming distance function
func TestMetasanoHammingDistance(t *testing.T) {
	const input1 = "this is a test"
	const input2 = "wokka wokka!!!"
	const expectedDistance = 37

	in1bytes, _ := hex.DecodeString(fmt.Sprintf("%x", input1))
	in2bytes, _ := hex.DecodeString(fmt.Sprintf("%x", input2))

	distance, _ := CalculateDistance(in1bytes, in2bytes)
	if distance != expectedDistance {
		t.Errorf("%v != %v", distance, expectedDistance)
	}
}

func TestMetasanoHammingDistanceWithEmptyArrays(t *testing.T) {
	const input1 = ""
	const input2 = ""
	const expectedDistance = 0

	in1bytes, _ := hex.DecodeString(fmt.Sprintf("%x", input1))
	in2bytes, _ := hex.DecodeString(fmt.Sprintf("%x", input2))

	distance, _ := CalculateDistance(in1bytes, in2bytes)
	if distance != expectedDistance {
		t.Errorf("%v != %v", distance, expectedDistance)
	}
}

func TestMetasanoHammingDistanceWithUnequalArrays(t *testing.T) {
	const input1 = "hello, world"
	const input2 = "goodbye, cruel world"
	const expectedDistance = -1

	in1bytes, _ := hex.DecodeString(fmt.Sprintf("%x", input1))
	in2bytes, _ := hex.DecodeString(fmt.Sprintf("%x", input2))

	distance, err := CalculateDistance(in1bytes, in2bytes)
	if err == nil {
		t.Fatal(distance, err)
	}
	if distance != expectedDistance {
		t.Errorf("%v != %v", distance, expectedDistance)
	}
}
