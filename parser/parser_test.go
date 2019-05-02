package parser

import (
	"fmt"
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
		t.Fatalf("program.Statements does not cotain 3 statements. got=%d",
			len(program.Statements))
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

/*
	parseLetStatement() のテスト
 */
func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q",
			s.TokenLiteral())
		return false
	}

	// 型アサーション : overrideした型の情報が欠落してしまうため, 実体の型を引数にしてチュックする.
	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s",
			name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

/*
	parseReturnStatement() のテスト
 */
func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
		`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
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

/*
	parseIdentifier() のテスト
 */
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()

	/*
		input を構文解析し, エラーがないか構文解析器を確認する.
	 */
	checkParserErrors(t, p)

	/*
	   input をパースした結果, *ast.Program ノードに含まれる文の数が 1つであることを確認する.
	*/
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	/*
		program.Statements に含まれる唯一の文が *ast.ExpressionStatement 型であることを確認する.
	 */
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	/*
		ExpressionStatement ノードの Expression が *ast.Identifier 型であることを確認して,
		Value と TokenLiteral が予測した出力結果と合っているか確認する.
	 */
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

/*
	parserIntegerLiteral() のテスト
 */
func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()

	/*
		input を構文解析し, エラーがないか構文解析器を確認する.
	*/
	checkParserErrors(t, p)

	/*
		input をパースした結果, *ast.Program ノードに含まれる文の数が 1つであることを確認する.
	*/
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	/*
		program.Statements に含まれる唯一の文が *ast.ExpressionStatement 型であることを確認する.
	*/
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	/*
		ExpressionStatement ノードの Expression が *ast.IntegerLiteral 型であることを確認して,
		Value と TokenLiteral が予測した出力結果と合っているか確認する.
	*/
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntgerLiteral. got=%d", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

/*
	parsePrefixExpression() のテスト
 */
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParserProgram()

		/*
			input を構文解析し, エラーがないか構文解析器を確認する.
		*/
		checkParserErrors(t, p)

		/*
			input をパースした結果, *ast.Program ノードに含まれる文の数が 1つであることを確認する.
		*/
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		/*
			program.Statements に含まれる唯一の文が *ast.ExpressionStatement 型であることを確認する.
		*/
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		/*
			ExpressionStatement ノードの Expression が *ast.PrefixExpression 型であることを確認して,
			Operator と Right が予測した出力結果と合っているか確認する.
		*/
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntergerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

/*
	parserIntegerLiteral() のテスト
 */
func testIntergerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}
