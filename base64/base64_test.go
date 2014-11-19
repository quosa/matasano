package base64

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// Test data from http://en.wikipedia.org/wiki/Base64
var base64testStrings = []struct {
	in       string
	expected string
}{
	{"Man", "TWFu"},
	{"Ma", "TWE="}, // Single padding
	{"M", "TQ=="},  // Double padding
	{"pleasure.", "cGxlYXN1cmUu"},
	{"leasure.", "bGVhc3VyZS4="},
	{"easure.", "ZWFzdXJlLg=="},
	{"asure.", "YXN1cmUu"},
	{"sure.", "c3VyZS4="},
	{"", ""}, // Special case, nothing to decode
}

/*
 * Test conversion from hexstring to base64
 */
func TestConvertSampleStringsToBase64(t *testing.T) {
	for _, test := range base64testStrings {
		inbytes, err := hex.DecodeString(fmt.Sprintf("%x", test.in))
		if err != nil {
			t.Fatal(err)
		}
		outbytes := ConvertToBase64(inbytes)
		outstring := fmt.Sprintf("%s", outbytes)
		if outstring != test.expected {
			t.Errorf("ConvertToBase64(%v) = %v, want %v", test.in, outstring, test.expected)
		}
	}
}

// from the matasano assignment
func TestConvertSampleHexStringToBase64(t *testing.T) {
	const in = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	const expected = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	inbytes, err := hex.DecodeString(in)
	if err != nil {
		t.Fatal(err)
	}
	outbytes := ConvertToBase64(inbytes)
	outstring := fmt.Sprintf("%s", outbytes)
	if outstring != expected {
		t.Errorf("ConvertToBase64(%v) = %v, want %v", in, outstring, expected)
	}
}

/*
 * Test conversion from hexstring to base64
 */
func TestConvertBase64ToString(t *testing.T) {
	for _, test := range base64testStrings {

		inbytes, err := hex.DecodeString(fmt.Sprintf("%x", test.expected))
		if err != nil {
			t.Fatal(err)
		}
		outbytes := ConvertFromBase64(inbytes)
		outstring := fmt.Sprintf("%s", outbytes)
		if outstring != test.in {
			t.Errorf("ConvertFromBase64(%v) = %v, want %v", test.expected, outstring, test.in)
		}
	}
}

// from the matasano assignment
func TestConvertSampleBase64ToHexString(t *testing.T) {
	const cipher = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	const expected = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	cipherbytes, err := hex.DecodeString(fmt.Sprintf("%x", cipher))
	if err != nil {
		t.Fatal(err)
	}

	outbytes := ConvertFromBase64(cipherbytes)

	outhex := fmt.Sprintf("%x", outbytes)
	if outhex != expected {
		t.Errorf("ConvertFromBase64(%v) = %v, want %v", cipher, outhex, expected)
	}
}

func TestConvertSinglePaddedBase64ToHexString(t *testing.T) {
	const cipher = "bGVhc3VyZS4="
	const expected = "leasure."
	cipherbytes, err := hex.DecodeString(fmt.Sprintf("%x", cipher))
	if err != nil {
		t.Fatal(err)
	}

	outbytes := ConvertFromBase64(cipherbytes)

	outstring := fmt.Sprintf("%s", outbytes)
	if outstring != expected {
		t.Errorf("ConvertFromBase64(%v) = %v (%v), want %v (%v)",
			cipher, outstring, len(outstring), expected, len(expected))
	}
}
