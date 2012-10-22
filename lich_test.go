package lich_test

import (
	"fmt"
	"github.com/earthboundkid/lich"
	"testing"
)

func ExampleData() {
	data := lich.Data("Hello, World!")

	fmt.Println(data)
	//Output: 13<Hello, World!>
}

func ExampleArray() {
	array := lich.Array{lich.Data("apple"), lich.Data("banana"), lich.Data("orange")}

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
		"greeting":       lich.Data([]byte{72, 101, 108, 108, 111}),
		"fruit":          lich.Array{lich.Data("apple"), lich.Data("banana"), lich.Data("orange")},
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
	data := lich.Data("")
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

func TestSimpleDataParsing(t *testing.T) {
	s := "5<hello>"
	element, err := lich.Parse(s)
	if err != nil || element != lich.Data("hello") {
		t.Fatalf("Parsed %q\nGot element:\t%#v, %s\nError:\t%#v", s, element, element, err)
	}

}

func TestSimpleArrayParsing(t *testing.T) {
	s := "26[5<apple>6<banana>6<orange>]"
	array := lich.Array{lich.Data("apple"), lich.Data("banana"), lich.Data("orange")}

	element, err := lich.Parse(s)
	if err != nil || element.String() != array.String() {
		t.Fatalf("Parsed %q\nGot element:\t%#v, %s\nError:\t%#v", s, element, element, err)
	}
}
