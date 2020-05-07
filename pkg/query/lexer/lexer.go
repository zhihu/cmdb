// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"github.com/zhihu/cmdb/pkg/query/token"
)

const (
	NoState    = -1
	NumStates  = 49
	NumSymbols = 101
)

type Lexer struct {
	src    []byte
	pos    int
	line   int
	column int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = new(token.Token)
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: '_'
1: '<'
2: '>'
3: '!'
4: '='
5: '<'
6: '='
7: '>'
8: '='
9: '='
10: '='
11: '='
12: '('
13: ')'
14: '!'
15: ','
16: 'n'
17: 'o'
18: 't'
19: 'i'
20: 'n'
21: 'i'
22: 'n'
23: 'e'
24: 'x'
25: 'i'
26: 's'
27: 't'
28: 's'
29: 'e'
30: 'x'
31: 'i'
32: 's'
33: 't'
34: 'n'
35: 'o'
36: 't'
37: 'e'
38: 'x'
39: 'i'
40: 's'
41: 't'
42: 'n'
43: 'o'
44: 't'
45: 'e'
46: 'x'
47: 'i'
48: 's'
49: 't'
50: 's'
51: 'A'
52: 'N'
53: 'D'
54: '&'
55: '&'
56: '#'
57: '$'
58: '%'
59: '&'
60: '''
61: '*'
62: '+'
63: '-'
64: '/'
65: '?'
66: '^'
67: '_'
68: '`'
69: '{'
70: '|'
71: '}'
72: '~'
73: '.'
74: '"'
75: '\'
76: ' '
77: '!'
78: '"'
79: '\'
80: '/'
81: 'b'
82: 'f'
83: 'n'
84: 'r'
85: 't'
86: 'u'
87: '\t'
88: '\n'
89: '\r'
90: ' '
91: 'a'-'z'
92: 'A'-'Z'
93: '0'-'9'
94: 'A'-'Z'
95: 'a'-'z'
96: '0'-'9'
97: \u0100-\U0010ffff
98: '#'-'['
99: ']'-\U0010ffff
100: .
*/
