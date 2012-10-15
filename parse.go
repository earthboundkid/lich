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
	r := rune(s[0])

	if !isdigit(r) {
		return nil, UnparseableError
	}

	start, stop := 0, 1
	for isdigit(rune(s[stop])) {
	}

	size, err := strconv.Atoi(s[start:stop])

	if err != nil {
		//I don't think this is reachable.
		return nil, UnparseableError
	}

	if rune(s[stop]) != '<' || rune(s[stop+size+1]) != '>' {
		return nil, UnparseableError
	}

	return Data(s[stop+1 : stop+size+1]), nil
}

func isdigit(r rune) bool {
	return r >= '0' && r <= '9'
}
