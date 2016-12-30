package lexer

import (
	"errors"
	"fmt"
	"github.com/owlci/gosonett/token"
	"unicode"
)

type LexerPosition struct {
	line     int
	lineChar int
}

func (lp *LexerPosition) NextLine() {
	lp.line++
	lp.lineChar = 0
}

func (lp *LexerPosition) NextChar() {
	lp.lineChar++
}

func (lp *LexerPosition) Format() string {
	return fmt.Sprintf("%w:%w", lp.line, lp.lineChar)
}

type Lexer struct {
	Source       string
	Tokens       []token.Token
	Position     LexerPosition // represents position of char with new lines, good for debugging.
	index        int           // hold the current index of currentChar within the whole input string
	sourceLength int
	reachedEnd   bool
}

const EOF = '\x00'

func New(source string) *Lexer {
	return &Lexer{
		Source:       source,
		Position:     LexerPosition{line: 0, lineChar: 0},
		index:        0,
		sourceLength: len(source),
		reachedEnd:   false,
	}
}

func (l *Lexer) willOverflow() bool {
	return l.index+1 >= l.sourceLength
}

// NOTE: This might need to rune, depending on what character set jsonet supports.
func (l *Lexer) CurrentChar() rune {
	if l.reachedEnd {
		return EOF
	}

	return rune(l.Source[l.index])
}

// The first time we reach the end we expect the calling code to handle it correctly, either
// by printing an error message for invalid source or by terminating the token. With this, we
// are explicitly including EOF to be a valid lexeme in the token string.
func (l *Lexer) invalidOverflow() bool {
	if l.willOverflow() {
		// We haven't reached the end yet
		if l.reachedEnd == false {
			l.reachedEnd = true
			return false
		}

		return true
	}

	return false
}

func (l *Lexer) NextChar() (rune, error) {
	if l.invalidOverflow() {
		return EOF, errors.New("Unhandled end of input looking for the next character")
	}

	char := l.CurrentChar()
	l.index++

	if char == '\n' {
		l.Position.NextLine()
	} else {
		l.Position.NextChar()
	}

	return char, nil
}

// Returns the next lookahead character without advancing the lexer
func (l *Lexer) Peek() (rune, error) {
	if l.invalidOverflow() {
		return EOF, errors.New("Unhandled end of input peeking the next character")
	}

	return rune(l.Source[l.index+1]), nil
}

// Advances through the whole string source and tokenizes every lexeme
func (l *Lexer) Lex() []token.Token {
	for r := l.Tokenize(); r.Type != token.EOF; r = l.Tokenize() {
	}

	return l.Tokens
}

// Returns the next valid token in the input stream
func (l *Lexer) Tokenize() token.Token {
	var tok token.Token

	l.eatWhitespace()
	char := l.CurrentChar()
	str := string(char)

	switch char {
	case EOF:
		tok = token.New(token.EOF, "(EOF)")
	case '{':
		tok = token.New(token.LBRACE, str)
	case '}':
		tok = token.New(token.RBRACE, str)
	case '[':
		tok = token.New(token.LBRACKET, str)
	case ']':
		tok = token.New(token.RBRACKET, str)
	case ',':
		tok = token.New(token.COMMA, str)
	case '.':
		tok = token.New(token.DOT, str)
	case '(':
		tok = token.New(token.LPAREN, str)
	case ')':
		tok = token.New(token.RPAREN, str)
	case ';':
		tok = token.New(token.SEMICOLON, str)
	case '!':
		tok = token.New(token.BANG, str)
	case '$':
		tok = token.New(token.DOLLAR, str)
	case ':':
		tok = token.New(token.COLON, str)
	case '~':
		tok = token.New(token.TILDE, str)
	case '+':
		tok = token.New(token.PLUS, str)
	case '-':
		tok = token.New(token.MINUS, str)
	case '&':
		tok = token.New(token.AMP, str)
	case '|':
		tok = token.New(token.PIPE, str)
	case '^':
		tok = token.New(token.CARET, str)
	case '=':
		tok = token.New(token.ASSIGN, str)
	case '<':
		tok = token.New(token.LANGLE, str)
	case '>':
		tok = token.New(token.RANGLE, str)
	case '*':
		tok = token.New(token.STAR, str)
	case '/':
		peekedChar, err := l.Peek()

		if err != nil {
			panic(err)
		}

		// Single-line comment
		if peekedChar == '/' {
			l.eatCurrentLine()
			return l.Tokenize()
		}

		// Multi-line comment
		if peekedChar == '*' {
			// TODO: Handle multi-line-comments
			// l.eatMultiLineComment()
			// return l.Tokenize()
		}

		// Must be a single token acting as an operator
		tok = token.New(token.SLASH, str)
	case '%':
		tok = token.New(token.PERC, str)
	case '#':
		l.eatCurrentLine()
		return l.Tokenize()
	case '"', '\'':
		// Whatever the opening char, we expect a closing char to match
		// but skip the first occurance since it starts the string
		l.NextChar()
		stringValue := l.eatUntil(char)
		tok = token.New(token.STRING, stringValue)
	// case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
	// token, _ := l.lexNumber()
	default:
		if isIdentifierFirst(char) {
			// NOTE: Error handling
			tok, _ = l.lexIdentifier()
		} else {
			// TODO: Use the LexerPosition struct to print out something nice here
			panic("Unknown lexing character")
		}
	}

	// Store the token
	l.Tokens = append(l.Tokens, tok)

	// End of token, advance to next byte
	l.NextChar()

	return tok
}

// Chews up insignificant whitespace up until the next potential token
func (l *Lexer) eatWhitespace() {
	// TODO: More idiomatic way to do this
	for unicode.IsSpace(rune(l.CurrentChar())) {
		l.NextChar()
	}
}

// TODO: This should panic if it doesn't find *untilChar* and reaches EOF
func (l *Lexer) eatUntil(untilChar rune) string {
	var eatenStr string

	for l.CurrentChar() != untilChar {
		char, err := l.NextChar()

		if err != nil {
			panic(err)
		}

		eatenStr = eatenStr + string(char)
	}

	return eatenStr
}

func (l *Lexer) eatUntilAfter(untilChar rune) string {
	eatenStr := l.eatUntil(untilChar)

	// Point to the byte after our untilChar
	l.NextChar()

	return eatenStr
}

func (l *Lexer) eatCurrentLine() {
	l.eatUntilAfter('\n')
}

// TODO...
func (l *Lexer) eatMultiLineComment() {
}

func (l *Lexer) lexIdentifier() (token.Token, error) {
	startIndex := l.index

	for isIdentifier(l.CurrentChar()) {
		char, err := l.NextChar()

		if err != nil {
			panic(err)
		}

		if char == EOF {
			break
		}
	}

	ident := l.Source[startIndex:l.index]

	// Backtrack one char to end on the last byte of the identifier/keyword
	l.index--

	// matchKeyword and return keyword token
	tokenType := token.GetKeywordKind(ident)

	return token.Token{Type: tokenType, Value: ident}, nil
}

// NOTE: Taken from here https://github.com/google/go-jsonnet/blob/master/lexer.go#L189
func isIdentifierFirst(r rune) bool {
	return unicode.IsUpper(r) || unicode.IsLower(r) || r == '_'
}

func isIdentifier(r rune) bool {
	return isIdentifierFirst(r) || unicode.IsNumber(r)
}
