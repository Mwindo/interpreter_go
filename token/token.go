package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT" // identifier
	INT   = "INT"

	// Operators
	ASSIGN          = "="
	PLUS            = "+"
	MINUS           = "-"
	BANG            = "!"
	ASTERISK        = "*"
	DOUBLE_ASTERISK = "**"
	CARET           = "^"
	SLASH           = "/"
	MODULUS         = "%"

	// Emojis
	EMOJI_TREE = "🌴"
	EMOJI_DINO = "🦖"

	EMOJI = "EMOJI"

	// Relations
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	//Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Scopes
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	STRING = "STRING"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
