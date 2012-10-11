package main

import (
	"bytes"
	"fmt"
	"sort"
)

//StrictMode controls whether or not DictFromMap panics when it encounters
//values of an unknown type.
var StrictMode = true

type Element interface {
	String() string
}

type Data string

func (data Data) String() string {
	return fmt.Sprintf("%d<%s>", len(data), string(data))
}

type Array []Element

func (array Array) String() string {
	var b bytes.Buffer

	for i := range array {
		b.WriteString(array[i].String())
	}

	return fmt.Sprintf("%d[%v]", b.Len(), b.String())
}

func ArrayFromStrings(strings ...string) Array {
	array := make(Array, 0, len(strings))
	for _, s := range strings {
		array = append(array, Data(s))
	}
	return array
}

type Dict map[Data]Element

func (d Dict) String() string {
	keys := make([]string, 0, len(d))
	for key := range d {
		keys = append(keys, string(key))
	}

	//Canonize order
	sort.Strings(keys)
	var b bytes.Buffer

	for _, key := range keys {
		de := Data(key)
		b.WriteString(de.String())
		b.WriteString(d[de].String())
	}

	return fmt.Sprintf("%d{%s}", b.Len(), b.String())
}

func DictFromMap(m map[string]interface{}) Dict {
	const panicmessage = "Cannot serialize unknown type.\nKey: %s\nValue: %v\nType: %T"

	d := make(Dict)
	for key := range m {
		switch value := m[key].(type) {
		case Data:
			d[Data(key)] = value
		case Array:
			d[Data(key)] = value
		case Dict:
			d[Data(key)] = value

		case string:
			d[Data(key)] = Data(value)
		case []byte:
			d[Data(key)] = Data(value)

		case []string:
			d[Data(key)] = ArrayFromStrings(value...)

		case map[string]interface{}:
			d[Data(key)] = DictFromMap(value)

		default:
			if StrictMode {
				panic(fmt.Sprintf(panicmessage, key, value, value))
			}
		}
	}
	return d
}

func main() {

	d1 := Dict{
		"selling points": ArrayFromStrings("simple", "general", "human-sympathetic"),
		"greeting":       Data([]byte{72, 101, 108, 108, 111}),
		"fruit":          Array{Data("apple"), Data("banana"), Data("orange")},
	}

	fmt.Println(d1)
	StrictMode = false
	d2 := map[string]interface{}{
		"greeting": map[string]interface{}{
			"English":  []byte{72, 101, 108, 108, 111},
			"Japanese": "Konnichiwa",
			"Dog":      9.0,
		},
		"fruit":          []string{"apple", "banana", "orange"},
		"selling points": ArrayFromStrings("simple", "general", "human-sympathetic"),
	}
	fmt.Println(DictFromMap(d2))
}
