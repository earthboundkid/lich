package lich

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type LexTokenType int

func (t LexTokenType) String() string {
	return []string{
		Data:       "Data",
		ArrayOpen:  "ArrayOpen",
		ArrayClose: "ArrayClose",
		DictOpen:   "DictOpen",
		DictClose:  "DictClose",
	}[int(t)]
}

// Lexer token enum
const (
	Data LexTokenType = iota
	ArrayOpen
	ArrayClose
	DictOpen
	DictClose
)

type LexToken struct {
	LexTokenType
	Data []byte
}

func (t LexToken) String() string {
	if t.LexTokenType == Data {
		return fmt.Sprintf("%v<%s>", t.LexTokenType, t.Data)
	}
	return t.LexTokenType.String()
}

// Errors introduced by Lex
var (
	ErrUnexpectedChar     = errors.New("Unexpected character")
	ErrBufferTooShort     = errors.New("Input too short")
	ErrMissingClosingChar = errors.New("Missing closing character")
)

func Lex(r io.Reader) ([]LexToken, error) {
	b := bufio.NewReader(r)
	return lexElement(b)
}

type readScanner interface {
	io.Reader
	io.ByteScanner
}

func lexElement(b readScanner) ([]LexToken, error) {
	// Spec says size will always fit into 20 bytes or less
	const maxSizeLength = 20

	sizeBuf := make([]byte, maxSizeLength)
	i := 0

	for i = range sizeBuf {
		c, err := b.ReadByte()
		if err != nil {
			return nil, err
		}

		if !isDigit(c) {
			b.UnreadByte()
			break
		}

		sizeBuf[i] = c
	}

	// Turn the size into an int
	size, err := strconv.Atoi(string(sizeBuf[:i]))
	if err != nil {
		return nil, ErrUnexpectedChar
	}

	c, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	// Bail out on bad data before making a possibly large buffer
	switch c {
	case '<', '[', '{':
	default:
		return nil, ErrUnexpectedChar
	}

	buf := make([]byte, size+1) // +1 for terminal character
	if _, err := io.ReadFull(b, buf); err != nil {
		return nil, err
	}

	switch c {
	case '<':
		return lexData(buf)
	case '[':
		return lexArrayOrDict(buf, true)
	case '{':
		return lexArrayOrDict(buf, false)
	}
	panic("unreachable")
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func lexData(buf []byte) ([]LexToken, error) {
	if len(buf) < 1 {
		return nil, ErrBufferTooShort
	}
	if buf[len(buf)-1] != '>' {
		return nil, ErrMissingClosingChar
	}
	tokens := []LexToken{{LexTokenType: Data, Data: buf[:len(buf)-1]}}
	return tokens, nil
}

func lexArrayOrDict(buf []byte, isArray bool) ([]LexToken, error) {
	var (
		closingChar = byte(']')
		openingType = ArrayOpen
		closingType = ArrayClose
	)
	if !isArray {
		closingChar = byte('}')
		openingType = DictOpen
		closingType = DictClose
	}

	// First check for consistency before parsing guts
	if buf[len(buf)-1] != closingChar {
		return nil, ErrMissingClosingChar
	}

	tokens := []LexToken{{LexTokenType: openingType}}

	b := bytes.NewReader(buf)

	for b.Len() > 1 {
		newTokens, err := lexElement(b)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, newTokens...)
	}

	tokens = append(tokens, LexToken{LexTokenType: closingType})
	return tokens, nil
}
