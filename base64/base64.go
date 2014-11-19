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

func ConvertToBase64(inbytes []byte) []byte {
	const caps = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lows = "abcdefghijklmnopqrstuvwxyz"
	const nums = "0123456789"
	const rest = "+/"
	const base64output = caps + lows + nums + rest

	outbytes := make([]byte, 0)

	padding := len(inbytes) % 3
	switch padding {
	case 1:
		inbytes = append(inbytes, byte(0x00), byte(0x00))
	case 2:
		inbytes = append(inbytes, byte(0x00))
	}
	for len(inbytes) >= 3 {
		h := inbytes[0] >> 2
		i := (inbytes[0]&byte(0x03))<<4 + inbytes[1]>>4
		j := (inbytes[1]&byte(0x0f))<<2 + inbytes[2]&byte(0xc0)>>6
		k := inbytes[2] & byte(0x3f)
		outbytes = append(outbytes, base64output[h], base64output[i], base64output[j], base64output[k])
		inbytes = inbytes[3:]
	}
	switch padding {
	case 1:
		outbytes[len(outbytes)-2] = 0x3d // '='
		outbytes[len(outbytes)-1] = 0x3d // '='
	case 2:
		outbytes[len(outbytes)-1] = 0x3d // '='
	}
	return outbytes
}

func ConvertFromBase64(base64bytes []byte) []byte {
	const caps = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lows = "abcdefghijklmnopqrstuvwxyz"
	const nums = "0123456789"
	const rest = "+/"
	const base64output = caps + lows + nums + rest

	base64input := make(map[rune]byte)
	for i, v := range base64output {
		base64input[v] = byte(i)
	} // NOTE: = -> \0 is implicit

	outbytes := make([]byte, 0)
	var chop int = 0
	base64len := len(base64bytes)
	if base64len > 2 && string(base64bytes[base64len-1]) == "=" {
		chop = 1
	}
	if base64len > 2 && string(base64bytes[base64len-2]) == "=" {
		chop = 2
	}

	for len(base64bytes) >= 4 {
		h, i, j, k := base64input[rune(base64bytes[0])],
			base64input[rune(base64bytes[1])],
			base64input[rune(base64bytes[2])],
			base64input[rune(base64bytes[3])]
		a := h<<2 + (i&byte(0x30))>>4
		b := i&byte(0x0f)<<4 + j&byte(0x3c)>>2
		c := j&byte(0x03)<<6 + k
		outbytes = append(outbytes, a, b, c)
		base64bytes = base64bytes[4:]
	}
	if chop > 0 && len(outbytes) > chop {
		outbytes = outbytes[:len(outbytes)-chop]
	}
	return outbytes
}
