package parser

import (
	"fmt"
	"github.com/WTBacon/goInterpreter/ast"
	"github.com/WTBacon/goInterpreter/lexer"
	"github.com/WTBacon/goInterpreter/token"
)

/*
	構文解析器（パーサー）を表す構造体型.
	l        	: 字句解析器インスタンスへのポインタ
	curToken 	: 現在調べているトークン
	peekToken 	: 次に調べるトークン
	errors		: 構文解析中のエラー
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
	字句解析器を受け取って構文解析器のインスタンスを生成する関数.
 */
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 2つのトークンを読み込む.
	// 1回目で, peekToken がセットされる.
	p.nextToken()
	// 2回目で, curToken　がセットされる.
	p.nextToken()
	return p
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
func (p *Parser) ParserProgram() *ast.Program {
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
		return nil
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
	任意のトークンタイプに遭遇するたびに, 対応する構文解析関数が呼ばれる.
	これらの関数は適切な式を構文解析し, その式を表現するASTノードを返す.
	トークンタイプごとに, 最大２つの構文解析関数が関連づけられる.
 */
type (
	prefixParseFn func() ast.Expression              // 前置構文解析関数（prefix parsing function）
	infixParseFn func(ast.Expression) ast.Expression // 中置構文解析関数（infix parsing function）
)

/*
	prefixParseFns マップと infixParseFns マップにエントリを追加するヘルパーメソッド.
 */
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
