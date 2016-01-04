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
	fmt.Stringer
	// Dummy private method to maintain package boundaries
	isElement()
}

type DataString string

func (d DataString) isElement() {}

func (data DataString) String() string {
	return fmt.Sprintf("%d<%s>", len(data), string(data))
}

type Array []Element

func (a Array) isElement() {}

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
		array = append(array, DataString(s))
	}
	return array
}

type Dict map[DataString]Element

func (d Dict) isElement() {}

func (d Dict) String() string {
	keys := make([]string, 0, len(d))
	for key := range d {
		keys = append(keys, string(key))
	}

	//Canonize order
	sort.Strings(keys)
	var b bytes.Buffer

	for _, key := range keys {
		de := DataString(key)
		b.WriteString(de.String())
		b.WriteString(d[de].String())
	}

	return fmt.Sprintf("%d{%s}", b.Len(), b.String())
}

func DictFromMap(m map[string]interface{}) (Dict, error) {
	d := make(Dict)
	for key := range m {
		switch value := m[key].(type) {
		case DataString:
			d[DataString(key)] = value
		case Array:
			d[DataString(key)] = value
		case Dict:
			d[DataString(key)] = value

		case string:
			d[DataString(key)] = DataString(value)
		case []byte:
			d[DataString(key)] = DataString(value)

		case []string:
			d[DataString(key)] = ArrayFromStrings(value...)

		case map[string]interface{}:
			subdict, err := DictFromMap(value)
			if err != nil {
				return Dict{}, err
			}
			d[DataString(key)] = subdict

		default:
			return Dict{}, UnknownTypeError{key, value}
		}
	}
	return d, nil
}
