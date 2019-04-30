package parser

import (
	"github.com/WTBacon/goInterpreter/ast"
	"github.com/WTBacon/goInterpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {

	/*
		インプットとなるソースコード.
	 */
	input := `
		let x = 5;
		let y = 10;
		let foobar = 8383838;
		`

	/*
		字句解析器にソースコードを与えて初期化.
	*/
	l := lexer.New(input)

	/*
		構文解析器に字句解析器を与えて初期化.
	*/
	p := New(l)

	/*
		構文解析器でパース.
	*/
	program := p.ParserProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not cotain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatements(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

/*
	構文解析器のエラーをチャックし, もしエラーがあればテストエラーとして, テストの実行を停止する.
 */
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
