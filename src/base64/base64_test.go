package base64

import (
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
		out := ConvertToBase64(fmt.Sprintf("%x", test.in))
		if out != test.expected {
			t.Errorf("ConvertToBase64(%v) = %v, want %v", test.in, out, test.expected)
		}
	}
}

// from the matasano assignment
func TestConvertSampleHexStringToBase64(t *testing.T) {
	const in = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	const expected = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	if out := ConvertToBase64(in); out != expected {
		t.Errorf("ConvertToBase64(%v) = %v, want %v", in, out, expected)
	}
}

/*
 * Test conversion from hexstring to base64
 */
func TestConvertBase64ToString(t *testing.T) {
	for _, test := range base64testStrings {
		out := ConvertFromBase64(test.expected)
		if out != test.in {
			t.Errorf("ConvertFromBase64(%v) = %v, want %v", test.expected, out, test.in)
		}
	}
}

// from the matasano assignment
func TestConvertSampleBase64ToHexString(t *testing.T) {
	const in = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	const expected = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	if out := ConvertFromBase64(in); fmt.Sprintf("%x", out) != expected {
		t.Errorf("ConvertFromBase64(%v) = %v, want %v", in, out, expected)
	}
}
