package lich

import (
	"bytes"
	"fmt"
	"sort"
)

//DictFromMap returns an UnknownTypeError 
//when it encounters values of an unknown type.
type UnknownTypeError struct {
	Key   string
	Value interface{}
}

func (u UnknownTypeError) Error() string {
	const panicmessage = "Cannot serialize unknown type.\nKey: %s\nValue: %v\nType: %T"
	return fmt.Sprintf(panicmessage, u.Key, u.Value, u.Value)
}

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

func DictFromMap(m map[string]interface{}) (Dict, error) {
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
			subdict, err := DictFromMap(value)
			if err != nil {
				return Dict{}, err
			}
			d[Data(key)] = subdict

		default:
			return Dict{}, UnknownTypeError{key, value}
		}
	}
	return d, nil
}
