/*
1. Convert hex to base64 and back.

The string:

49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d

should produce:

SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t

Now use this code everywhere for the rest of the exercises. Here's a
simple rule of thumb:

Always operate on raw bytes, never on encoded strings. Only use hex
and base64 for pretty-printing.
*/
package base64

import "encoding/hex"

func ConvertToBase64(hexstring string) string {
	const caps = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lows = "abcdefghijklmnopqrstuvwxyz"
	const nums = "0123456789"
	const rest = "+/"
	const base64output = caps + lows + nums + rest

	var out string = ""

	data, err := hex.DecodeString(hexstring)
	if err != nil {
		panic(err)
	}

	padding := len(data) % 3
	switch padding {
	case 1:
		data = append(data, byte(0x00), byte(0x00))
	case 2:
		data = append(data, byte(0x00))
	}
	for len(data) >= 3 {
		h := data[0] >> 2
		i := (data[0]&byte(0x03))<<4 + data[1]>>4
		j := (data[1]&byte(0x0f))<<2 + data[2]&byte(0xc0)>>6
		k := data[2] & byte(0x3f)
		out += string(base64output[h])
		out += string(base64output[i])
		out += string(base64output[j])
		out += string(base64output[k])
		data = data[3:]
	}
	switch padding {
	case 1:
		out = out[:len(out)-2] + "=="
	case 2:
		out = out[:len(out)-1] + "="
	}
	return out
}

func ConvertFromBase64(base64 string) string {
	const caps = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lows = "abcdefghijklmnopqrstuvwxyz"
	const nums = "0123456789"
	const rest = "+/"
	const base64output = caps + lows + nums + rest

	base64input := make(map[rune]byte)
	for i, v := range base64output {
		base64input[v] = byte(i)
	} // NOTE: = -> \0 is implicit

	var out string = ""
	var chop int = 0
	if len(base64) > 0 && string(base64[len(base64)-1]) == "=" {
		chop = 1
	}
	if len(base64) > 0 && string(base64[len(base64)-2]) == "=" {
		chop = 2
	}
	for len(base64) >= 4 {
		h, i, j, k := base64input[rune(base64[0])],
			base64input[rune(base64[1])],
			base64input[rune(base64[2])],
			base64input[rune(base64[3])]
		a := h<<2 + (i&byte(0x30))>>4
		b := i&byte(0x0f)<<4 + j&byte(0x3c)>>2
		c := j&byte(0x03)<<6 + k
		out += string(string(a) + string(b) + string(c))
		base64 = base64[4:]
	}
	if chop > 0 && len(out) > chop {
		out = out[:len(out)-chop]
	}
	return out
}
