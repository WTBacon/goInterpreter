package lexer

import "github.com/WTBacon/goInterpreter/token"

/*
	字句解析器（レキサー）を表す構造体型.
	input			: ソースコード
	position 		: 常に最後に読んだ位置を示す（chの位置を示すインデクス）
	readPosition 	: 次に読み込む位置を示す
	ch         		: 現在検査中の文字
}
 */
type Lexer struct {
	input        string // ソースコード
	position     int    // 常に最後に読んだ位置を示す（chの位置を示すインデクス）
	readPosition int    // 次に読み込む位置を示す
	ch           byte   // 現在検査中の文字
}

/*
	ソースコード（input） から 字句解析器（Lexer 型の構造体）を生成.
	readChar() で初期化.
 */
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

/*
	ソースコードの次の一文字（readPosition）を読んで, 現在位置（position）を進める.
	「ch = 0」は「まだ何も読み込んでいない」もしくは「ファイルの終わり」を表す.
	TODO: Bacon で Unicode と絵文字をサポートする.
 */
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

/*
	現在検査中の文字（ch） に一致する Bacon の Token を返す.
	Token を返す前に, 入力のポインタを返す.
 */
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {

	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupTokenType(tok.Literal)
			return tok
		} else if isDisit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

/*
	ch が 識別子 / キーワードの一部であれば,
	読み終えるまでポインタを進めて, 読み込んだ識別子 / キーワードの文字列を返す.
 */
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

/*
	ch が整数であれば, 読み終えるまでポインタを進めて, 読み込んだ整数を文字列で返す.
 */
func (l *Lexer) readNumber() string {
	position := l.position
	for isDisit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

/*
	ソースコードを正確に TokenType にパースするために, 先読みするためのヘルパーメソッド.
	先読みした文字だけを返す.
 */
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

/*
	スペースやタブ, 改行を読み飛ばすためのヘルパーメソッド.
 */
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

/*
	予期しない文字が来た時に, token.ILLEGAL トークンとして扱うためのヘルパー関数.
 */
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

/*
	与えられた文字が, 英字もしくは"_"か判定するヘルパー関数.
 */
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

/*
	与えられた文字が, 数字か判定するヘルパー関数.
 */
func isDisit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
