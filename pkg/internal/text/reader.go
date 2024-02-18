package text

import (
	"bufio"
	"io"
	"strings"
)

type FilteredRune rune

const (
	SoftHyphen      = FilteredRune('\u00AD')
	ZeroWidthSpace  = FilteredRune('\u200B')
	Comma           = FilteredRune(',')
	Dot             = FilteredRune('.')
	RoundBraceLeft  = FilteredRune('(')
	RoundBraceRight = FilteredRune(')')
	QuestionMark    = FilteredRune('?')
	ExclamationMark = FilteredRune('!')
)

type CustomRuneReader struct {
	r      *bufio.Reader
	filter map[rune]bool
}

func GetDefaultFilteredRunes() []rune {
	return []rune{
		rune(Dot),
		rune(Comma),
		rune(SoftHyphen),
		rune(ZeroWidthSpace),
		rune(RoundBraceLeft),
		rune(RoundBraceRight),
		rune(QuestionMark),
		rune(ExclamationMark),
	}
}

var nonWords = GetDefaultNonWords()

func GetDefaultNonWords() map[string]struct{} {
	m := make(map[string]struct {
	})

	m["-"] = struct{}{}

	return m
}

func NewCustomRuneReader(r *bufio.Reader, filterRunes ...rune) *CustomRuneReader {
	filter := make(map[rune]bool)
	for _, r := range filterRunes {
		filter[r] = true
	}
	return &CustomRuneReader{
		r:      r,
		filter: filter,
	}
}

func Words(r io.Reader, filterRunes ...rune) ([]string, error) {

	filter := make(map[rune]bool)
	for _, r := range filterRunes {
		filter[r] = true
	}

	words := make([]string, 0)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {

		trimmed := strings.TrimFunc(scanner.Text(), func(r rune) bool {
			return filter[r]
		})

		if _, ok := nonWords[trimmed]; ok {
			continue
		}

		words = append(words, trimmed)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

func (cr *CustomRuneReader) ReadRune() (rune, int, error) {
	r, size, err := cr.r.ReadRune()
	if err != nil {
		return r, size, err
	}

	if !cr.filter[r] {
		return r, size, nil
	} else {
		return cr.ReadRune()
	}

}
