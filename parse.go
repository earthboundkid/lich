package lich

import (
	"errors"
	"strings"

	"github.com/carlmjohnson/lich/lex"
)

// Errors introduced during decoding
var (
	ErrUnexpectedClose = errors.New("Got unexpected close token")
	ErrMissingClose    = errors.New("Did not receive a close token")
	ErrBadKeyType      = errors.New("Lich dict key must be of type data")
)

func Decode(s string) (Element, error) {
	sc := lex.NewScanner(strings.NewReader(s))
	el, err := decodeElement(sc)
	if err != nil {
		return nil, err
	}
	if sc.Next() {
		panic("Huh?")
	}
	if err = sc.Error(); err != nil {
		return nil, err
	}
	return el, err
}

func decodeElement(s *lex.Scanner) (Element, error) {
	if !s.Next() {
		return nil, ErrUnexpectedClose
	}

	switch s.Token {
	case lex.DataToken:
		return DataString(s.Data), nil
	case lex.ArrayOpen:
		return decodeArray(s)
	case lex.ArrayClose:
		return nil, ErrUnexpectedClose
	case lex.DictOpen:
		return decodeDict(s)
	case lex.DictClose:
		return nil, ErrUnexpectedClose
	}
	panic("Unknown lexer token type")
}

func decodeArray(s *lex.Scanner) (Element, error) {
	var a Array
	for s.Next() {
		switch s.Token {
		case lex.DataToken:
			a = append(a, DataString(s.Data))
		case lex.ArrayOpen:
			el, err := decodeArray(s)
			if err != nil {
				return nil, err
			}
			a = append(a, el)

		case lex.ArrayClose:
			return a, nil
		}
	}
	return nil, ErrMissingClose
}

func decodeDict(s *lex.Scanner) (Element, error) {
	d := make(Dict)
	for s.Next() {
		switch s.Token {
		case lex.DataToken:
			key := DataString(s.Data)
			el, err := decodeElement(s)
			if err != nil {
				return nil, err
			}
			d[key] = el
		case lex.DictClose:
			return d, nil
		default:
			return nil, ErrBadKeyType
		}
	}
	return nil, ErrMissingClose
}
