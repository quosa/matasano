/*
2. Fixed XOR

Write a function that takes two equal-length buffers and produces
their XOR sum.

The string:

1c0111001f010100061a024b53535009181c

... after hex decoding, when xor'd against:

686974207468652062756c6c277320657965

... should produce:

746865206b696420646f6e277420706c6179
*/

package xorsum

import "errors"

func FixedXORSum(data1, data2 []byte) ([]byte, error) {
	len1 := len(data1)
	if len1 != len(data2) {
		return nil, errors.New("buffers must have equal length!")
	}
	result := make([]byte, len1)
	for i := 0; i < len1; i++ {
		result[i] = data1[i] ^ data2[i]
	}
	return result, nil
}
