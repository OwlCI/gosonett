package lexer

import (
	"github.com/owlci/gosonett/token"
	"testing"
)

type TokenMatcher struct {
	expectedType  token.TokenType
	expectedValue string
}

func runTokenMatches(t *testing.T, source string, tests []TokenMatcher) {
	lexer := New(source)
	tokens := lexer.Lex()
	testsLength := len(tests)
	tokensLength := len(tokens)

	if testsLength != tokensLength {
		t.Fatalf("Wrong token array length: expected=%d, got=%d", tokensLength, testsLength)
	}

	for i, tok := range tokens {
		tm := tests[i]

		if tok.Type != tm.expectedType {
			t.Fatalf("Wrong token type: expected=%q, got=%q", tm.expectedType, tok.Type)
		}

		if tok.Value != tm.expectedValue {
			t.Fatalf("Wrong token value: expected=%q, got=%q", tm.expectedValue, tok.Value)
		}
	}
}

func TestSymbols(t *testing.T) {
	source := "{}[],.();"

	tests := []TokenMatcher{
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.COMMA, ","},
		{token.DOT, "."},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestOperators(t *testing.T) {
	source := "!$:~+-&|^=<>*/%"

	tests := []TokenMatcher{
		{token.BANG, "!"},
		{token.DOLLAR, "$"},
		{token.COLON, ":"},
		{token.TILDE, "~"},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.AMP, "&"},
		{token.PIPE, "|"},
		{token.CARET, "^"},
		{token.ASSIGN, "="},
		{token.LANGLE, "<"},
		{token.RANGLE, ">"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.PERC, "%"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestWhitepace(t *testing.T) {
	source := "! =        %"

	tests := []TokenMatcher{
		{token.BANG, "!"},
		{token.ASSIGN, "="},
		{token.PERC, "%"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestComments(t *testing.T) {
	source := `
! # Inline Comment !!!!!
= // Inline comment ===
!
`

	tests := []TokenMatcher{
		{token.BANG, "!"},
		{token.ASSIGN, "="},
		{token.BANG, "!"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestKeywords(t *testing.T) {
	source := `
assert
error
if
then
else
true
false
for
function
import
importstr
tailstrict
in
local
null
self
super
`

	tests := []TokenMatcher{
		{token.ASSERT, "assert"},
		{token.ERROR, "error"},
		{token.IF, "if"},
		{token.THEN, "then"},
		{token.ELSE, "else"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.FOR, "for"},
		{token.FUNCTION, "function"},
		{token.IMPORT, "import"},
		{token.IMPORTSTR, "importstr"},
		{token.TAILSTRICT, "tailstrict"},
		{token.IN, "in"},
		{token.LOCAL, "local"},
		{token.NULL, "null"},
		{token.SELF, "self"},
		{token.SUPER, "super"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestIdentifiers(t *testing.T) {
	source := "_testThis UpperW1thNum"

	tests := []TokenMatcher{
		{token.IDENT, "_testThis"},
		{token.IDENT, "UpperW1thNum"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestString(t *testing.T) {
	source := `"double" 'single'`

	tests := []TokenMatcher{
		{token.STRING, "double"},
		{token.STRING, "single"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}

func TestSnippet(t *testing.T) {
	source := `
// Jsonnet Example
{
	person1: {
		name: "Alice",
		welcome: "Hello " + self.name + "!",
	},
	person2: self.person1 { name: "Bob" },
}
`
	tests := []TokenMatcher{
		{token.LBRACE, "{"},
		{token.IDENT, "person1"},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.STRING, "Alice"},
		{token.COMMA, ","},
		{token.IDENT, "welcome"},
		{token.COLON, ":"},
		{token.STRING, "Hello "},
		{token.PLUS, "+"},
		{token.SELF, "self"},
		{token.DOT, "."},
		{token.IDENT, "name"},
		{token.PLUS, "+"},
		{token.STRING, "!"},
		{token.COMMA, ","},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.IDENT, "person2"},
		{token.COLON, ":"},
		{token.SELF, "self"},
		{token.DOT, "."},
		{token.IDENT, "person1"},
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.STRING, "Bob"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.RBRACE, "}"},
		{token.EOF, "(EOF)"},
	}

	runTokenMatches(t, source, tests)
}
