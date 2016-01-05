package lex_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/carlmjohnson/lich/lex"
)

func TestScannerErrors(t *testing.T) {
	var tests = []struct {
		input string
		err   bool
	}{
		{"", true},
		{"0<>", false},
		{"0[]", false},
		{"0{}", false},
		{"0{}x", true},
		{"1<>", true},
		{"1<a>", false},
		{"26{8<greeting>11<hello world>}", false},
		{"26[5<apple>6<banana>6<orange>]", false},
		{"126{14<selling points>40[6<simple>7<general>17<human-sympathetic>]" +
			"8<greeting>11<hello world>5<fruit>26[5<apple>6<banana>6<orange>]}",
			false},
	}
	for _, test := range tests {
		fmt.Println(test.input)
		r := strings.NewReader(test.input)
		s := lex.NewScanner(r)
		for s.Next() {
		}

		if (s.Error() != nil) != test.err {
			t.Errorf("lex.FromString(%q).error = %v", test.input, s.Error())
		}
	}
}
