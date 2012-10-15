package lich

import (
	"fmt"
	"strconv"
)

type UnparseableError struct {
	Parsestring string
	Location    int
	Problem     string
}

const errformat = "Couldn't parse string %q...\nProblem at index %d was %q."

func (u UnparseableError) Error() string {
	return fmt.Sprintf(errformat, u.Parsestring[:10], u.Location, u.Problem)
}

func Parse(s string) (Element, error) {
	return topLevel(s, 0, len(s))
}

func isdigit(r uint8) bool {
	return r >= '0' && r <= '9'
}

func topLevel(s string, start, stop int) (Element, error) {
	if len(s) < 1 {
		return nil, UnparseableError{s, 0, "Empty string!"}
	}

	current := start
	for isdigit(s[current]) {
		current++
	}

	size, err := strconv.Atoi(s[start:current])

	if err != nil {
		return nil, UnparseableError{s, current, "Non-digit start"}
	}

	//If this doesn't match, the reported size is screwed up.
	//Doing this check helps make sure we don't try to read too far.
	if current+size+2 != stop {
		return nil, UnparseableError{s, current, "Data payload is too short"}
	}

	switch s[current] {
	case '<':
		if s[stop-1] != '>' {
			return nil, UnparseableError{s, stop - 1, "No matching >"}
		}
		return Data(s[current+1 : stop-1]), nil

	case '[':
		if s[stop-1] != ']' {
			return nil, UnparseableError{s, stop - 1, "No matching ]"}
		}

		return getArray(s, current+1, stop-1), nil

	}
	return nil, UnparseableError{s, current, "Invalid separator"}
}

func getArray(s string, start, stop int) Element {
	if s[stop] != ']' {
		return nil
	}
	return Data(s[start:stop])
}
