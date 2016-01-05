package lich_test

import (
	"fmt"
	"testing"

	"github.com/carlmjohnson/lich"
)

func ExampleData() {
	data := lich.DataString("Hello, World!")

	fmt.Println(data)
	//Output: 13<Hello, World!>
}

func ExampleArray() {
	array := lich.Array{lich.DataString("apple"), lich.DataString("banana"), lich.DataString("orange")}

	fmt.Println(array)
	//Output: 26[5<apple>6<banana>6<orange>]
}

func ExampleArrayFromStrings() {
	array := lich.ArrayFromStrings("simple", "general", "human-sympathetic")

	fmt.Println(array)
	//Output: 40[6<simple>7<general>17<human-sympathetic>]

}

func ExampleDict() {
	d1 := lich.Dict{
		"selling points": lich.ArrayFromStrings("simple", "general", "human-sympathetic"),
		"greeting":       lich.DataString([]byte{72, 101, 108, 108, 111}),
		"fruit":          lich.Array{lich.DataString("apple"), lich.DataString("banana"), lich.DataString("orange")},
	}

	fmt.Println(d1)
	//Output: 119{5<fruit>26[5<apple>6<banana>6<orange>]8<greeting>5<Hello>14<selling points>40[6<simple>7<general>17<human-sympathetic>]}
}

func ExampleDictFromMap() {
	d2, _ := lich.DictFromMap(map[string]interface{}{
		"greeting": map[string]interface{}{
			"English":  []byte{72, 101, 108, 108, 111},
			"Japanese": "Konnichiwa",
		},
		"fruit":          []string{"apple", "banana", "orange"},
		"selling points": lich.ArrayFromStrings("simple", "general", "human-sympathetic"),
	})
	fmt.Println(d2)
	//Output: 158{5<fruit>26[5<apple>6<banana>6<orange>]8<greeting>43{7<English>5<Hello>8<Japanese>10<Konnichiwa>}14<selling points>40[6<simple>7<general>17<human-sympathetic>]}

}

func TestEmptyData(t *testing.T) {
	data := lich.DataString("")
	str := data.String()
	if str != "0<>" {
		t.Fatal(data, str)
	}
}

func TestEmptyArray(t *testing.T) {
	array := lich.Array{}
	str := array.String()
	if str != "0[]" {
		t.Fatal(array, str)
	}
}

func TestEmptyDict(t *testing.T) {
	d := lich.Dict{}
	str := d.String()
	if str != "0{}" {
		t.Fatal(d, str)
	}
}

func TestInvalidMap(t *testing.T) {
	d2, err := lich.DictFromMap(map[string]interface{}{
		"greeting": map[string]interface{}{
			"English":  []byte{72, 101, 108, 108, 111},
			"Japanese": "Konnichiwa",
			"Dog":      9.0,
		},
		"fruit":          []string{"apple", "banana", "orange"},
		"selling points": lich.ArrayFromStrings("simple", "general", "human-sympathetic"),
	})
	if (err != lich.UnknownTypeError{"Dog", 9.0}) {
		t.Fatal(d2, err)
	}
}

var decodeTests = []struct {
	input string
	err   bool
}{
	{"", true},
	{"0<>", false},
	{"0[]", false},
	{"4[1<a>]", false},
	{"7[4[1<a>]]", false},
	{"26{8<greeting>11<hello world>}", false},
	{"26[5<apple>6<banana>6<orange>]", false},
	{"126{5<fruit>26[5<apple>6<banana>6<orange>]8<greeting>11<hello world>" +
		"14<selling points>40[6<simple>7<general>17<human-sympathetic>]}",
		false},
}

func TestDecodeErrors(t *testing.T) {
	for _, test := range decodeTests {
		_, err := lich.Decode(test.input)
		if (err != nil) != test.err {
			t.Errorf("Decode(%v).err = %v", test.input, err)
		}
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	for _, test := range decodeTests {
		if test.err {
			continue
		}
		el, err := lich.Decode(test.input)
		if err != nil {
			t.Errorf("Decode(%v).err = %v", test.input, err)
		}
		if el == nil {
			t.Errorf("Decode(%v) = %v", test.input, el)
		}
		if el.String() != test.input {
			t.Errorf("Decode(%v) = %v", test.input, el)
		}
	}
}
