package toki

import (
	"testing"
)

func TestScanner(t *testing.T) {
	const (
		NUMBER Token = iota + 1
		PLUS
		STRING
	)

	input := "1  + 2+3 + happy birthday  "
	t.Logf("input: %s", input)
	scanner := NewScanner(
		[]Def{
			{Token: NUMBER, Pattern: "[0-9]+"},
			{Token: PLUS, Pattern: `\+`},
			{Token: STRING, Pattern: "[a-z]+"},
		})
	s := scanner.Scan(input)
	expected := []Token{
		NUMBER,
		PLUS,
		NUMBER,
		PLUS,
		NUMBER,
		PLUS,
		STRING,
		STRING,
		EOF,
	}
	for _, e := range expected {
		r := s.Next()
		if e != r.Token {
			t.Fatalf("expected %v, got %v", e, r.Token)
		} else {
			t.Log(r)
		}
	}
}
