package toki

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode/utf8"
)

var regexSpace = regexp.MustCompile(`^[\t ]+`)

type Token uint32

const (
	EOF Token = 1<<32 - 1 - iota
	Error
)

type Position struct {
	Line   int
	Column int
}

func (p *Position) move(s []byte) {
	p.Line += bytes.Count(s, []byte{'\n'})
	last := bytes.LastIndex(s, []byte{'\n'})
	if last != -1 {
		p.Column = 1
	}
	p.Column += utf8.RuneCount(s[last+1:])
}

type Def struct {
	Token   Token
	Pattern string
	regexp  *regexp.Regexp
}

type Result struct {
	Token Token
	Value []byte
	Pos   Position
}

func (r *Result) String() string {
	return fmt.Sprintf("Line: %d, Column: %d, %s", r.Pos.Line, r.Pos.Column, r.Value)
}

type Scanner struct {
	space *regexp.Regexp
	def   []Def
}

func NewScanner(def []Def) *Scanner {
	for i := range def {
		def[i].regexp = regexp.MustCompile("^" + def[i].Pattern)
	}
	return &Scanner{
		space: regexSpace,
		def:   def,
	}
}

func (s *Scanner) Scan(input string) *Session {
	return &Session{
		scanner: s,
		pos:     Position{1, 1},
		input:   []byte(input),
	}
}

type Session struct {
	scanner *Scanner
	pos     Position
	input   []byte
}

func (s *Session) skip() {
	result := s.scanner.space.Find(s.input)
	if result == nil {
		return
	}
	s.pos.move(result)
	s.input = bytes.TrimPrefix(s.input, result)
}

func (s *Session) scan() *Result {
	s.skip()
	if len(s.input) == 0 {
		return &Result{Token: EOF, Pos: s.pos}
	}
	for _, r := range s.scanner.def {
		result := r.regexp.Find(s.input)
		if result == nil {
			continue
		}
		return &Result{Token: r.Token, Value: result, Pos: s.pos}
	}
	return &Result{Token: Error, Pos: s.pos}
}

func (s *Session) Peek() *Result {
	return s.scan()
}

func (s *Session) Next() *Result {
	t := s.scan()
	if t.Token == Error || t.Token == EOF {
		return t
	}
	s.input = bytes.TrimPrefix(s.input, t.Value)
	s.pos.move(t.Value)
	return t
}
