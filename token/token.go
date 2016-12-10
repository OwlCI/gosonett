package token

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

func New(tokenType TokenType, val string) Token {
	return Token{Type: tokenType, Value: val}
}

const (
	// Special Lexemes
	EOF   = "EOF"
	IDENT = "IDENT"

	// TODO: Symbols
	// {}[],.();
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	COMMA     = ","
	DOT       = "."
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"

	// Operators
	// !$:~+-&|^=<>*/%
	MINUS  = "-"
	PLUS   = "+"
	BANG   = "!"
	TILDE  = "~"
	DOLLAR = "$"
	COLON  = ":"
	AMP    = "&"
	PIPE   = "|"
	CARET  = "^"
	ASSIGN = "="
	LANGLE = "<"
	RANGLE = ">"
	STAR   = "*"
	SLASH  = "/"
	PERC   = "%"

	// Keywords
	ASSERT     = "assert"
	ERROR      = "error"
	IF         = "if"
	THEN       = "then"
	ELSE       = "else"
	TRUE       = "true"
	FALSE      = "false"
	FOR        = "for"
	FUNCTION   = "function"
	IMPORT     = "import"
	IMPORTSTR  = "importstr"
	TAILSTRICT = "tailstrict"
	IN         = "in"
	LOCAL      = "local"
	NULL       = "null"
	SELF       = "self"
	SUPER      = "super"
)

var keywords = map[string]TokenType{
	"assert":     ASSERT,
	"error":      ERROR,
	"if":         IF,
	"then":       THEN,
	"else":       ELSE,
	"true":       TRUE,
	"false":      FALSE,
	"for":        FOR,
	"function":   FUNCTION,
	"import":     IMPORT,
	"importstr":  IMPORTSTR,
	"tailstrict": TAILSTRICT,
	"in":         IN,
	"local":      LOCAL,
	"null":       NULL,
	"self":       SELF,
	"super":      SUPER,
}

func GetKeywordKind(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
