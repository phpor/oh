// Released under an MIT license. See LICENSE.

package parser

import (
	. "github.com/michaelmacinnis/oh/pkg/cell"
	"github.com/michaelmacinnis/oh/pkg/common"
	"github.com/michaelmacinnis/oh/pkg/system"
	"io"
	"os"
	"strings"
)

type parser struct {
	deref func(string, uintptr) Cell
}

type scanner struct {
	*parser
	error    func(file string, line int, text string)
	f        *os.File
	filename string
	input    common.ReadStringer
	process  func(Cell, string, int, string) (Cell, bool)

	line []rune

	cursor   int
	lineno   int
	previous rune
	start    int
	state    int
	token    rune

	finished bool
}

const (
	ssStart = iota
	ssAmpersand
	ssBang
	ssBangDouble
	ssBangDoubleEscape
	ssBangGreater
	ssColon
	ssComment
	ssDoubleQuoted
	ssDoubleQuotedEscape
	ssGreater
	ssLess
	ssPipe
	ssSingleQuoted
	ssSymbol
)

func (s *scanner) Lex(lval *yySymType) (token int) {
	var operator = map[string]string{
		"!>":  "_redirect_stderr_",
		"!>>": "_append_stderr_",
		"!|":  "_pipe_stderr_",
		"!|+": "_channel_stderr_",
		"&":   "spawn",
		"&&":  "and",
		"<":   "_redirect_stdin_",
		"<(":  "_substitute_stdout_",
		">":   "_redirect_stdout_",
		">(":  "_substitute_stdin_",
		">>":  "_append_stdout_",
		"|":   "_pipe_stdout_",
		"|+":  "_channel_stdout_",
		"||":  "or",
	}

	defer func() {
		exists := false

		v := string(s.line[s.start:s.cursor])

		switch s.token {
		case SYMBOL:
			if strings.ContainsAny(v, "{}") {
				if v == "{" || v == "}" {
					s.token = s.line[s.start]
				} else {
					s.token = BRACE_EXPANSION
				}
				token = int(s.token)
			}
			lval.s = v

		case BACKGROUND, ORF, ANDF, PIPE, REDIRECT, SUBSTITUTE:
			lval.s, exists = operator[v]
			if exists {
				break
			}
			lval.s = v

		default:
			lval.s = v
		}

		s.state = ssStart
		s.previous = s.token
		s.token = 0
	}()

	retries := 0

main:
	for s.token == 0 {
		if s.cursor >= len(s.line) {
			if s.finished {
				return 0
			}

			line, err := s.input.ReadString('\n')
			if err == nil {
				retries = 0
			} else if err == common.CtrlCPressed {
				s.start = 0
				s.token = CTRLC
				break
			} else if s.f != nil && retries < 1 && err != io.EOF {
				if system.ResetForegroundGroup(s.f) {
					retries++
					goto main
				}
			}

			s.lineno++

			runes := []rune(line)
			last := len(runes) - 2
			if last >= 0 && runes[last] == '\r' {
				runes = append(runes[0:last], '\n')
				last--
			}

			if last >= 0 && runes[last] == '\\' {
				runes = runes[0:last]
			}

			if err != nil {
				runes = append(runes, '\n')
				s.finished = true
			}

			if s.start < s.cursor-1 {
				s.line = append(s.line[s.start:s.cursor], runes...)
				s.cursor -= s.start
			} else {
				s.cursor = 0
				s.line = runes
			}
			s.start = 0
			s.token = 0
		}

		switch s.state {
		case ssStart:
			s.start = s.cursor

			switch s.line[s.cursor] {
			default:
				s.state = ssSymbol
				continue main
			case '\n', '%', '(', ')', ';', '@', '`', '}':
				s.token = s.line[s.start]
			case '\t', ' ':
				s.state = ssStart
			case '!':
				s.state = ssBang
			case '"':
				s.state = ssDoubleQuoted
			case '#':
				s.state = ssComment
			case '&':
				s.state = ssAmpersand
			case '\'':
				s.state = ssSingleQuoted
			case ':':
				s.state = ssColon
			case '<':
				s.state = ssLess
			case '>':
				s.state = ssGreater
			case '|':
				s.state = ssPipe
			}

		case ssAmpersand:
			switch s.line[s.cursor] {
			case '&':
				s.token = ANDF
			default:
				s.token = BACKGROUND
				continue main
			}

		case ssBang:
			switch s.line[s.cursor] {
			case '"':
				s.state = ssBangDouble
			case '>':
				s.state = ssBangGreater
			case '|':
				s.state = ssPipe
			default:
				s.state = ssSymbol
				continue main
			}

		case ssBangGreater:
			s.token = REDIRECT
			if s.line[s.cursor] != '>' {
				continue main
			}

		case ssColon:
			switch s.line[s.cursor] {
			case ':':
				s.token = CONS
			default:
				s.token = ':'
				continue main
			}

		case ssComment:
			for s.line[s.cursor] != '\n' ||
				s.line[s.cursor-1] == '\\' {
				s.cursor++

				if s.cursor >= len(s.line) {
					continue main
				}
			}
			s.cursor--
			s.state = ssStart

		case ssBangDouble, ssBangDoubleEscape,
			ssDoubleQuoted, ssDoubleQuotedEscape:
			for s.cursor < len(s.line) {
				if s.state == ssBangDoubleEscape {
					s.state = ssBangDouble
				} else if s.state == ssDoubleQuotedEscape {
					s.state = ssDoubleQuoted
				} else if s.line[s.cursor] == '"' {
					break
				} else if s.line[s.cursor] == '\\' {
					if s.state == ssBangDouble {
						s.state = ssBangDoubleEscape
					} else {
						s.state = ssDoubleQuotedEscape
					}
				}
				s.cursor++
			}
			if s.cursor >= len(s.line) {
				continue main
			}
			if s.state == ssBangDouble {
				s.token = BANG_DOUBLE
			} else {
				s.token = DOUBLE_QUOTED
			}

		case ssGreater:
			s.token = REDIRECT
			if s.line[s.cursor] == '(' {
				s.token = SUBSTITUTE
			} else if s.line[s.cursor] != '>' {
				continue main
			}

		case ssLess:
			s.token = REDIRECT
			if s.line[s.cursor] == '(' {
				s.token = SUBSTITUTE
			} else {
				continue main
			}

		case ssPipe:
			switch s.line[s.cursor] {
			case '+':
				s.token = PIPE
			case '|':
				s.token = ORF
			default:
				s.token = PIPE
				continue main
			}

		case ssSingleQuoted:
			for s.cursor < len(s.line) && s.line[s.cursor] != '\'' {
				s.cursor++
			}
			if s.cursor >= len(s.line) {
				if s.line[s.cursor-1] == '\n' {
					s.line = append(s.line[0:s.cursor-1], []rune("\\n")...)
				}
				continue main
			}
			s.token = SINGLE_QUOTED

		case ssSymbol:
			switch s.line[s.cursor] {
			case '\n', '%', '&', '\'', '(', ')', ';',
				'<', '@', '`', '|',
				'\t', ' ', '"', '#', ':', '>':
				s.token = SYMBOL
				continue main
			}

		}
		s.cursor++

		if s.token == '\n' {
			switch s.previous {
			case ORF, ANDF, PIPE, REDIRECT:
				s.token = 0
			}
		}
	}

	return int(s.token)
}

func (s *scanner) Error(msg string) {
	s.error(s.filename, s.lineno, msg)
}

func New(deref func(string, uintptr) Cell) *parser {
	return &parser{deref}
}

func (p *parser) Parse(
	input common.ReadStringer,
	error func(file string, line int, text string), f *os.File,
	filename string, process func(Cell, string, int, string) (Cell, bool),
) bool {

	s := new(scanner)

	s.parser = p
	s.error = error
	s.f = f
	s.filename = filename
	s.input = input
	s.process = process

	rval := 1
	for rval > 0 {
		s.line = []rune("")
		s.cursor = 0
		s.previous = 0
		s.start = 0
		s.token = 0

		s.finished = false

		s.state = ssStart

		rval = yyParse(s)
	}

	return rval == 0
}

//go:generate go tool yacc -o grammar.go grammar.y
//go:generate sed -i.save -f grammar.sed grammar.go
//go:generate go fmt grammar.go
//go:generate rm -f y.output grammar.go.save
