package lich

import (
	"errors" //Why would anyone want to import more errors into their code???
	"strconv"
)

var UnparseableError = errors.New("Couldn't parse string.")

func Parse(s string) (Element, error) {
	if len(s) < 1 {
		return nil, UnparseableError
	}

	el := getElement(s, 0, len(s))

	if el == nil {
		return nil, UnparseableError
	}

	return el, nil
}

func isdigit(r uint8) bool {
	return r >= '0' && r <= '9'
}

func getElement(s string, start, stop int) Element {
	r := s[start]

	if !isdigit(r) {
		return nil
	}

	current := start + 1
	for isdigit(s[current]) {
		current++
	}

	size, _ := strconv.Atoi(s[start:current])

	//If this doesn't match, the reported size is screwed up.
	//Doing this check helps make sure we don't try to read too far.
	if current+size+2 <= stop {
		return nil
	}

	switch s[current] {
	case '<':
		return getData(s, current+1, stop)
	}
	return nil
}

func getData(s string, start, stop int) Element {
	if s[stop-1] != '>' {
		return nil
	}
	return Data(s[start : stop-1])
}
