package parser

import (
	"fmt"
	"github.com/WTBacon/goInterpreter/ast"
	"github.com/WTBacon/goInterpreter/lexer"
	"github.com/WTBacon/goInterpreter/token"
	"strconv"
)

/*
	構文解析器（パーサー）を表す構造体型.
	l        		: 字句解析器インスタンスへのポインタ
	curToken 		: 現在調べているトークン
	peekToken 		: 次に調べるトークン
	errors			: 構文解析中のエラー
	prefixParseFns	: 前置構文解析関数のマップ
	infixParseFns 	: 中置構文解析関数のマップ
}
 */
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

/*
	任意のトークンタイプに遭遇するたびに, 対応する構文解析関数が呼ばれる.
	これらの関数は適切な式を構文解析し, 式を表現するASTノード（Expressionノード）を返す.
	トークンタイプごとに, 最大２つの構文解析関数が関連づけられる.
 */
type (
	prefixParseFn func() ast.Expression              // 前置構文解析関数（prefix parsing function）
	infixParseFn func(ast.Expression) ast.Expression // 中置構文解析関数（infix parsing function）
)

/*
	prefixParseFns マップにエントリを追加するヘルパーメソッド.
 */
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

/*
	infixParseFns マップにエントリを追加するヘルパーメソッド.
 */
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

/*
	字句解析器を受け取って構文解析器のインスタンスを生成する関数.
 */
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 前置構文解析関数の初期化
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// IDENT トークンは, Identifier ノードにパースする.
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	// INT トークンは, IntegerLiteral ノードにパースする.
	p.registerPrefix(token.INT, p.parserIntegerLiteral)
	// Prefix となるトークンは, PrefixExpression ノードにパースする.
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// 真偽値トークンは, Boolean ノードにパースする.
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// LPAREN トークンは, グループ化された式としてパースする.
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// 中置構文解析関数の初期化
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// 以下のトークンは, InfixExpression ノードにパースする.
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	// 2つのトークンを読み込む.
	// 1回目で, peekToken がセットされる.
	p.nextToken()
	// 2回目で, curToken　がセットされる.
	p.nextToken()
	return p
}

/*
	現在のトークンを Identifier ノードにパースするメソッド.
 */
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

/*
	構文解析中のエラーを返すヘルパーメソッド.
 */
func (p *Parser) Errors() []string {
	return p.errors
}

/*
	curToken と peekToken を進める Parser のヘルパーメソッド.
 */
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

/*
	パースして抽象構文木を出力するメソッド.
 */
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

/*
	文をパースするメソッド.
	現在検査しているトークンを見て, どの文に一致するか判定する.
 */
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

/*
	let 文をパースするメソッド.
	LetStatement インスタンスを生成して, let 文が終了するまでトークンのポインタを進める.
 */
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// 識別子の名前を格納.
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

/*
	return 文をパースするメソッド.
	ReturnStatement インスタンスを生成して, return 文が終了するまでトークンのポインタを進める.
 */
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: セミコンに遭遇するまで式を読み飛ばしてしまっている.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

/*
	後続のトークンの型をチェックして, その方が正しい場合に限って nextToken を呼ぶアサーション関数.
 */
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

/*
	peekToken に期待していないトークンが来た時にエラー処理をするメソッド.
 */
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

/*
	現在調べるトークンが引数のトークンに一致するか判定するヘルパーメソッド.
 */
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

/*
	次に調べるトークンが引数のトークンに一致するか判定するヘルパーメソッド.
 */
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

/*
	式をパースするメソッド.
	式を表すトークンの前置解析関数をマップから入手して, 構文解析して Expression ノードを返す.
 */
func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFns[p.curToken.Type]
	// curToken が前置演算子の場合のみパースを継続する.
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// peekToken の左結合力（peekPrecedence()）が, curTokenの右結合力（引数の precedence）より高ければ,
	// これまで構文解析したもの（leftExp）は, 次の演算子に吸収される（infix(leftExp)）.
	// グループ化された式をパースするとき, 演算子やリテラルが続く限り左に結合していく.（")"は LOWEST になる）
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

/*
	式文をパースするメソッド.
	ExpressionStatement インスタンスを生成して, 式文が終了するまでトークンのポインタを進める.
 */
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

/*
	peekToken のトークンタイプに対応している優先順位を返すメソッド.
	対応している優先順位がなければ, デフォルト値で LOWEST を返す.
 */
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

/*
	curToken のトークンタイプに対応している優先順位を返すメソッド.
	対応している優先順位がなければ, デフォルト値で LOWEST を返す.
 */
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

/*
	式をパースした時に, 予期しない prefix に遭遇したら（prefixParseFnsになかったら）,
	パーサーに Error を追加するメソッド.
 */
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESSGREATER  // > または <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X または !X
	CALL         // myFunction(X)
)

/*
	整数リテラルをパースするメソッド.
	IntegerLiteral インスタンスを生成して, 現在参照しているトークンのリテラルを Int 型にパースして,
	IntegerLiteral インスタンスに入れて IntegerLiteral ノードを返す.
 */
func (p *Parser) parserIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))

	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

/*
	前置演算子を含む式をパースするメソッド.
	PrefixExpression インスタンスを生成して, 前演算子を含む式をPrefixExpression ノードにパースして返す.
 */
func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))

	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	// PREFIX より優先度が高いトークン（LPAREN）に遭遇しない限り,
	// 後続のトークンが expression.Right としてパースされる.
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

/*
	トークンタイプの優先順位マップ : トークンタイプとその優先順位を関連づける.
 */
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,      // =
	token.NOT_EQ:   EQUALS,      // !=
	token.LT:       LESSGREATER, // <
	token.GT:       LESSGREATER, // >
	token.PLUS:     SUM,         // +
	token.MINUS:    SUM,         // -
	token.SLASH:    PRODUCT,     // /
	token.ASTERISK: PRODUCT,     // *
	token.LPAREN:   CALL,        // )
}

/*
	中置演算子を含む式をパースするメソッド.
	InfixExpression インスタンスを生成して, 中演算子を含む式をInfixExpression ノードにパースして返す.
	curToken が 中置演算子のときに呼ばれる.
	ex. <expression>(Left) <infix operator>(curToken) <expression>(Right)
 */
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// 中置演算子の優先順位を保持する.
	precedence := p.curPrecedence()
	// curToken を前進させる.
	p.nextToken()
	// 中置演算子の優先順位を保持しながら, Right の式をパースする.
	expression.Right = p.parseExpression(precedence)

	return expression

}

/*
	真偽値リテラルをパースするメソッド.
	Boolean インスタンスを生成して, Booleanノードにパースして返す.
 */
func (p *Parser) parseBoolean() ast.Expression {
	defer untrace(trace("parseBoolean"))
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

/*
	グループ化された式をパースするメソッド.
	curToken が LPAREN トークン（"("）のときに呼び出され,
	最初の parseExpression において RPAREN トークン（")"）の優先順位（LOWEST）が参照されるまでパースする.
*/
func (p *Parser) parseGroupedExpression() ast.Expression {
	defer untrace(trace("parseGroupedExpression"))
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
