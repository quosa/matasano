package xorsum

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestMetasanoFixedXORSum(t *testing.T) {
	const hexstring1 = "1c0111001f010100061a024b53535009181c"
	const hexstring2 = "686974207468652062756c6c277320657965"
	const expected = "746865206b696420646f6e277420706c6179"

	data1, _ := hex.DecodeString(hexstring1)
	data2, _ := hex.DecodeString(hexstring2)

	out, err := FixedXORSum(data1, data2)
	if err != nil {
		t.Fatal(err)
	}

	hexout := hex.EncodeToString(out)
	if hexout != expected {
		t.Errorf("FixedXORSum(%v, %v) = %v, want %v",
			hexstring1, hexstring2, hexout, expected)
	}
}

func TestEmptyXORSum(t *testing.T) {
	data1 := make([]byte, 0)
	data2 := data1
	expected := data1
	out, err := FixedXORSum(data1, data2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, expected) {
		t.Errorf("FixedXORSum([], []) = %v, want []", hex.EncodeToString(out))
	}
}

func TestDifferentLengthInputForXORSum(t *testing.T) {
	data1 := []byte{1, 2}
	data2 := []byte{1, 2, 3, 4}
	out, err := FixedXORSum(data1, data2)
	if err == nil || out != nil {
		t.Fatal(out, err)
	}
	errorString := string(err.Error())
	if errorString != "buffers must have equal length!" {
		t.Errorf("FixedXORSum([2], [4]) raised %v, want "+
			"\"buffers must have equal length!\"", errorString)
	}

}
