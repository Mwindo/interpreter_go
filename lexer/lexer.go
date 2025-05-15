package lexer

import (
	"interpreter/token"
	"unicode/utf8"
)

type Lexer struct {
	input         string
	position      int  // current position in input
	readPosition  int  // next position in input
	currentSymbol rune // current symbol under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readSymbol()
	return l
}

func (l *Lexer) readSymbol() {
	if l.readPosition >= len(l.input) {
		l.currentSymbol = 0
	} else {
		// Decode the next rune from the input
		var size int
		l.currentSymbol, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.position = l.readPosition
		l.readPosition += size
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currentSymbol {
	case '=':
		if l.peekSymbol() == '=' {
			symbol := l.currentSymbol
			l.readSymbol()
			tok = token.Token{Type: token.EQ, Literal: string(symbol) + string(l.currentSymbol)}
		} else {
			tok = newToken(token.ASSIGN, l.currentSymbol)
		}
	case '+':
		tok = newToken(token.PLUS, l.currentSymbol)
	case '-':
		tok = newToken(token.MINUS, l.currentSymbol)
	case '!':
		if l.peekSymbol() == '=' {
			symbol := l.currentSymbol
			l.readSymbol()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(symbol) + string(l.currentSymbol)}
		} else {
			tok = newToken(token.BANG, l.currentSymbol)
		}
	case '/':
		tok = newToken(token.SLASH, l.currentSymbol)
	case '*':
		if l.peekSymbol() == '=' {
			symbol := l.currentSymbol
			l.readSymbol()
			tok = token.Token{Type: token.DOUBLE_ASTERISK, Literal: string(symbol) + string(l.currentSymbol)}
		} else {
			tok = newToken(token.ASTERISK, l.currentSymbol)
		}
	case '^':
		tok = newToken(token.CARET, l.currentSymbol)
	case '<':
		tok = newToken(token.LT, l.currentSymbol)
	case '>':
		tok = newToken(token.GT, l.currentSymbol)
	case ';':
		tok = newToken(token.SEMICOLON, l.currentSymbol)
	case ',':
		tok = newToken(token.COMMA, l.currentSymbol)
	case '(':
		tok = newToken(token.LPAREN, l.currentSymbol)
	case ')':
		tok = newToken(token.RPAREN, l.currentSymbol)
	case '{':
		tok = newToken(token.LBRACE, l.currentSymbol)
	case '}':
		tok = newToken(token.RBRACE, l.currentSymbol)
	case '🌴':
		tok = newToken(token.EMOJI_TREE, l.currentSymbol)
	case '🦖':
		tok = newToken(token.EMOJI_TREE, l.currentSymbol)
	case '%':
		tok = newToken(token.MODULUS, l.currentSymbol)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.currentSymbol) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.currentSymbol) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else if isEmoji(l.currentSymbol) {
			tok.Type = token.EMOJI
			tok.Literal = string(l.currentSymbol)
		} else {
			tok = newToken(token.ILLEGAL, l.currentSymbol)
		}
	}

	l.readSymbol()
	return tok
}

func newToken(tokenType token.TokenType, symbol rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(symbol)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.currentSymbol) {
		l.readSymbol()
	}
	return l.input[position:l.position]
}

func isLetter(symbol rune) bool {
	return 'a' <= symbol && symbol <= 'z' || 'A' <= symbol && symbol <= 'Z' || symbol == '_'
}

func isEmoji(symbol rune) bool {
	// Check against some known Unicode ranges for emoji
	return (symbol >= 0x1F600 && symbol <= 0x1F64F) || // Emoticons
		(symbol >= 0x1F300 && symbol <= 0x1F5FF) || // Miscellaneous Symbols and Pictographs
		(symbol >= 0x1F680 && symbol <= 0x1F6FF) || // Transport and Map Symbols
		(symbol >= 0x2600 && symbol <= 0x26FF) || // Miscellaneous Symbols
		(symbol >= 0x2700 && symbol <= 0x27BF) || // Dingbats
		(symbol >= 0x1F1E6 && symbol <= 0x1F1FF) // Regional Indicator Symbols (flags)
}

func (l *Lexer) skipWhitespace() {
	for l.currentSymbol == ' ' || l.currentSymbol == '\t' || l.currentSymbol == '\n' || l.currentSymbol == '\r' {
		l.readSymbol()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.currentSymbol) {
		l.readSymbol()
	}
	return l.input[position:l.position]
}

func isDigit(symbol rune) bool {
	return '0' <= symbol && symbol <= '9'
}

func (l *Lexer) peekSymbol() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readSymbol()
		if l.currentSymbol == '"' || l.currentSymbol == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
