package lex

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

// Token is an enum for the Decoder's lex tokens
type Token int

// Auto-stringer:
//go:generate stringer -type=Token

// Lexer token enum
const (
	Nil Token = iota
	DataToken
	ArrayOpen
	ArrayClose
	DictOpen
	DictClose
)

// Errors introduced by lexer
var (
	ErrUnexpectedChar     = errors.New("Unexpected character")
	ErrUnreadableSize     = errors.New("Couldn't read size of element")
	ErrMissingClosingChar = errors.New("Missing closing character")
)

// Scanner scans for the next lich token in the provided io.Reader on
// each call to Next().
type Scanner struct {
	next  func() bool
	err   error
	Token Token
	Data  []byte
}

// NewScanner takes an io.Reader and returns a Scanner.
func NewScanner(r io.Reader) *Scanner {
	s := &Scanner{}
	b := bufio.NewReader(r)

	s.next = func() bool {
		// We should get an EOF the next time this is called.
		s.next = func() bool {
			_, err := b.ReadByte()
			if err == io.EOF {
				return s.setError(nil)
			}
			if err != nil {
				return s.setError(err)
			}
			return s.setError(ErrUnexpectedChar)
		}
		return s.lexElement(b)
	}
	return s
}

// Next scans the underlying buffer for a token. If it finds a token,
// it sets Token (and optionally Data) and returns true. If it
// encounters an error or finishes scanning, it returns false.
func (s *Scanner) Next() bool {
	return s.next()
}

func (s *Scanner) setError(err error) bool {
	s.next = func() bool { return false }
	s.err = err
	return false
}

func (s *Scanner) Error() error {
	return s.err
}

func (s *Scanner) lexElement(b *bufio.Reader) bool {

	// Spec says size will always fit into 20 bytes or less
	const maxSizeLength = 20

	sizeBuf := make([]byte, maxSizeLength)
	i := 0

	for i = range sizeBuf {
		c, err := b.ReadByte()
		if err != nil {
			return s.setError(err)
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
		return s.setError(ErrUnreadableSize)
	}

	c, err := b.ReadByte()
	if err != nil {
		return s.setError(err)
	}

	switch c {
	case '<':
		return s.lexData(b, size)
	case '[':
		return s.lexArrayOrDict(b, size, true)
	case '{':
		return s.lexArrayOrDict(b, size, false)
	}

	return s.setError(ErrUnexpectedChar)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (s *Scanner) lexData(b *bufio.Reader, size int64) bool {
	data := make([]byte, int(size))
	_, err := io.ReadFull(b, data)
	if err != nil {
		return s.setError(err)
	}

	c, err := b.ReadByte()
	if err != nil {
		return s.setError(err)
	}
	if c != '>' {
		return s.setError(ErrMissingClosingChar)
	}

	s.Token = DataToken
	s.Data = data
	return true
}

func (s *Scanner) lexArrayOrDict(b *bufio.Reader, size int64, isArray bool) bool {
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

	lb := bufio.NewReader(io.LimitReader(b, size))

	origin := s.next
	s.next = func() bool {
		_, err := lb.ReadByte()
		if err == io.EOF {
			c, err := b.ReadByte()
			if err != nil {
				return s.setError(err)
			}
			if c != closingChar {
				return s.setError(ErrMissingClosingChar)
			}
			s.next = origin
			s.Token = closingType
			s.Data = nil
			return true
		}
		if err != nil {
			return s.setError(err)
		}
		lb.UnreadByte()
		return s.lexElement(lb)
	}

	s.Token = openingType
	s.Data = nil
	return true
}
