package repeatingkeyxorcipher

import (
	"../base64"
	"../singlecharxorcipher"
	//"encoding/base64"
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"os"
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

func TestRepeatingKeyXORCipherForBreaking(t *testing.T) {
	// base64-encoded repeating-key XOR
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

	base64bytes := base64.ConvertToBase64(out)
	fmt.Println(fmt.Sprintf("%s", base64bytes))

	// Starting to backtrack
	cipherBytes := base64.ConvertFromBase64(base64bytes)

	// test the keysize guessing
	for keysize := 2; keysize < 20; keysize++ {
		sliceA := cipherBytes[:keysize]
		sliceB := cipherBytes[keysize : 2*keysize]
		distance, err := CalculateDistance(sliceA, sliceB)
		if err != nil || distance < 0 {
			t.Fatal(distance, err)
		}
		normalizedDistance := float32(distance) / float32(keysize)
		fmt.Printf("keysize %v gives dist %v\n", keysize, normalizedDistance)
	}
	// --> gives high number for the correct value
	// lock the correct keysize
	leastKeysize := 3

	transpose := make(map[int][]byte)
	for i, v := range cipherBytes {
		transpose[i%leastKeysize] = append(transpose[i%leastKeysize], v)
	}

	repeatingKey := make([]byte, leastKeysize)
	for i := 0; i < leastKeysize; i++ {
		key, score, err := singlecharxorcipher.BreakSingleCharXORCipher(transpose[i])
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("got score %v for transpose %v\n", score, i)
		repeatingKey[i] = key
	}

	fmt.Printf("guessed key %v and correct is %v\n", repeatingKey, keybytes)
	if !bytes.Equal(repeatingKey, keybytes) {
		t.Fatal("key mismatch")
	}

	plainBytes := EncodeRepeatingKeyXORCipher(cipherBytes, repeatingKey)
	plainText := fmt.Sprintf("%s", plainBytes)
	//fmt.Printf("%v vs %v", input, plainText)
	if input != plainText {
		t.Fatal("plain text mismatch")
	}

	t.Fatal("BOOM")

}

// 6. Break repeating-key XOR
func TestMetasanoBreakRepeatingKeyXORCipher(t *testing.T) {
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

	//fmt.Printf("Got %v of base64 encoded ciphertext\n", len(cipherText))
	//fmt.Printf("<%v>\n", cipherText)

	base64CipherBytes, err := hex.DecodeString(fmt.Sprintf("%x", base64CipherText))
	if err != nil {
		t.Fatal(err)
	}
	cipherBytes := base64.ConvertFromBase64(base64CipherBytes)
	//fmt.Printf("<%x>\n", cipherBytes)

	leastDistance := float32(math.MaxFloat32)
	leastKeysize := -1

	//fmt.Printf("len cipherBytes = %v \n", len(cipherBytes))

	for keysize := 2; keysize < 40; keysize++ {
		sliceA := cipherBytes[:keysize]
		sliceB := cipherBytes[keysize : 2*keysize]
		distance, err := CalculateDistance(sliceA, sliceB)
		if err != nil || distance < 0 {
			t.Fatal(distance, err)
		}
		normalizedDistance := float32(distance) / float32(keysize)
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

	//e. Now that you probably know the KEYSIZE: break the ciphertext into
	//blocks of KEYSIZE length.

	for leastKeysize := 2; leastKeysize < 40; leastKeysize++ {

		//f. Now transpose the blocks: make a block that is the first byte of
		//every block, and a block that is the second byte of every block, and
		//so on.
		transpose := make(map[int][]byte)
		for i, v := range cipherBytes {
			transpose[i%leastKeysize] = append(transpose[i%leastKeysize], v)
		}
		//fmt.Printf("transpose is %v\n", transpose)

		//g. Solve each block as if it was single-character XOR.
		//You already have code to do this.
		repeatingKey := make([]byte, leastKeysize)
		for i := 0; i < leastKeysize; i++ {
			key, score, err := singlecharxorcipher.BreakSingleCharXORCipher(transpose[i])
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("got score %v for transpose %v\n", score, i)
			repeatingKey[i] = key
		}
		// TRY: 5(1,2), 3(2), 2(2,5), 13(2,5384614)

		//e. For each block, the single-byte XOR key that produces the best
		//looking histogram is the repeating-key XOR key byte for that
		//block. Put them together and you have the key.

		out := EncodeRepeatingKeyXORCipher(cipherBytes, repeatingKey)

		//fmt.Printf("%v\n%v", out, fmt.Sprintf("%s", out))
		fmt.Printf("\n\n\n\n\nleastKeysize %v\n", leastKeysize)
		fmt.Printf("%v\n", fmt.Sprintf("%s", out))
		//fmt.Println(hex.EncodeToString(out))
	}
	t.Fatal("BOOM")

	//if hex.EncodeToString(out) != expected {
	//	t.Errorf("%v (len %v)", out, len(out))
	//}

	//if strings.TrimSpace(topResult) != "Now that the party is jumping" {
	//	t.Errorf("%v (len %v)", topResult, len(topResult))
	//}

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