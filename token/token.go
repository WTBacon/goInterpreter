package token

/*
	Bacon 言語におけるトークンタイプを示す型.
 */
type TokenType string

/*
	Bacon 言語におけるトークンを表す構造体.
 */
type Token struct {
	Type    TokenType
	Literal string
}

/*
	Bacon 言語におけるキーワード.
 */
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

/*
	Token の Literal から, その Token の TokenTypeを探す関数.
 */
func LookupTokenType(literal string) TokenType {
	if tok, ok := keywords[literal]; ok {
		return tok
	}
	return IDENT
}

/*
	Bacon 言語におけるトークンタイプの種類.
 */
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 + リテラル
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1234567

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)
