package lex

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TokenType is an enum for the lexer's tokens
type TokenType int

// Auto-stringer:
//go:generate stringer -type=TokenType

// Lexer token enum
const (
	DataToken TokenType = iota
	ArrayOpen
	ArrayClose
	DictOpen
	DictClose
)

// Token has a TokenType. Tokens of type DataToken also have Data []byte.
type Token struct {
	TokenType
	Data []byte
}

func (t Token) String() string {
	if t.TokenType == DataToken {
		return fmt.Sprintf("%v<%s>", t.TokenType, t.Data)
	}
	return t.TokenType.String()
}

// Errors introduced by lexer
var (
	ErrUnexpectedChar     = errors.New("Unexpected character")
	ErrUnreadableSize     = errors.New("Couldn't read size of element")
	ErrMissingClosingChar = errors.New("Missing closing character")
)

// FromReader takes an io.Reader and returns lexer tokens or an error.
func FromReader(r io.Reader) ([]Token, error) {
	b := bufio.NewReader(r)
	return lexElement(b)
}

func lexElement(b *bufio.Reader) ([]Token, error) {
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

func lexData(b *bufio.Reader, size int64) ([]Token, error) {
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

	tokens := []Token{{TokenType: DataToken, Data: data}}
	return tokens, nil
}

func lexArrayOrDict(b *bufio.Reader, size int64, isArray bool) ([]Token, error) {
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
	tokens := []Token{{TokenType: openingType}}

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
	tokens = append(tokens, Token{TokenType: closingType})
	return tokens, nil
}

// FromString is a convenience method for lexing strings.
func FromString(s string) ([]Token, error) {
	r := strings.NewReader(s)
	return FromReader(r)
}

// FromBytes is a convenience method for lexing a slice of bytes.
func FromBytes(b []byte) ([]Token, error) {
	r := bytes.NewReader(b)
	return FromReader(r)
}
