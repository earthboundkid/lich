package lich

import (
	"bufio"
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
	ErrUnreadableSize     = errors.New("Couldn't read size of element")
	ErrMissingClosingChar = errors.New("Missing closing character")
)

func Lex(r io.Reader) ([]LexToken, error) {
	b := bufio.NewReader(r)
	return lexElement(b)
}

func lexElement(b *bufio.Reader) ([]LexToken, error) {
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
	size, err := strconv.ParseInt(string(sizeBuf[:i]), 10, 64)
	if err != nil {
		return nil, ErrUnreadableSize
	}

	c, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	switch c {
	case '<':
		return lexData(b, size)
	case '[':
		return lexArrayOrDict(b, size, true)
	case '{':
		return lexArrayOrDict(b, size, false)
	}

	return nil, ErrUnexpectedChar
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func lexData(b *bufio.Reader, size int64) ([]LexToken, error) {
	data := make([]byte, int(size))
	_, err := io.ReadFull(b, data)
	if err != nil {
		return nil, err
	}

	c, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if c != '>' {
		return nil, ErrMissingClosingChar
	}

	tokens := []LexToken{{LexTokenType: Data, Data: data}}
	return tokens, nil
}

func lexArrayOrDict(b *bufio.Reader, size int64, isArray bool) ([]LexToken, error) {
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
	tokens := []LexToken{{LexTokenType: openingType}}

	lb := bufio.NewReader(io.LimitReader(b, size))

	for {
		_, err := lb.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lb.UnreadByte()
		newTokens, err := lexElement(lb)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, newTokens...)
	}

	c, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if c != closingChar {
		return nil, ErrMissingClosingChar
	}
	tokens = append(tokens, LexToken{LexTokenType: closingType})
	return tokens, nil
}
