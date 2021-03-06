package ast

import (
	"bytes"
	"github.com/WTBacon/goInterpreter/token"
	"strings"
)

/*
	抽象構文木（Abstract sytax tree; AST）: ソースコードを入力とした構文解析器の出力.
	以下は, 再帰下降構文解析器のという AST の実装.
*/

/*
	AST の全てのノードは Node interface を実装しなければならない.
	つまり, TokenLiteral() メソッドを override する.
	TokenLiteral()	: ノードに関連づけられているトークンのリテラル値を返す.
	String()		: デバッグ時に AST ノードの情報を表示したり, 他のASTノードと比較したりする.
*/
type Node interface {
	TokenLiteral() string
	String() string
}

/*
	Statement（文）を表すノード.
	statementNode() : ダミーメソッド. コンパイルの段階で弾かせるため実装は持たなくて良い.
 */
type Statement interface {
	Node
	statementNode()
}

/*
	Expression（式）を表すノード.
	expressionNode() : ダミーメソッド. コンパイルの段階で弾かせるため実装は持たなくて良い.
 */
type Expression interface {
	Node
	expressionNode()
}

/*
	構文解析器が生成する全ての AST のルートノードを示す構造体型.
	一続きの文の集まりを格納するため, Statement インターフェースを実装する AST ノードのスライス.
 */
type Program struct {
	Statements []Statement
}

/*
	バッファを作成して, それぞれの文の String() メソッドの戻り値をバッファに書き込み, 文字列として返す.
 */
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

/*
	ルートノードのトークンのリテラルを返すメソッド.
 */
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

/*
	let 文を表す構造体型.（ex. let <identifier> = <expression>;）
	Toke 	: let 文を示すトークン
	Name	: 識別子の名前
	Value	: 値を生成する式
 */
type LetStatement struct {
	Token token.Token // token.LET トークン
	Name  *Identifier
	Value Expression
}

/*
	Node インターフェースと Statement インターフェースを override.
 */
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

/*
	識別子を表す構造体型.
	Token 	: 識別子を示すトークン
	Value	: 識別子の値
 */
type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

/*
	Node インターフェースと Expression インターフェースを override.
	なぜ Expresison なのかというと, 識別子は値を生成するから.（ex. let x = valueIdentifier;）
	ノードの種類を少なく保ち, 式として識別子を表現する.
 */
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

/*
	return 文を表す構造体型.（ex. return <expression>;）
	Toke	 	: let 文を示すトークン
	ReturnValue	: 値を返す式
 */
type ReturnStatement struct {
	Token       token.Token // 'return' トークン
	ReturnValue Expression
}

/*
	Node インターフェースと Statement インターフェースを override.
 */
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

/*
	式文を表す構造体型. (ex. x + 10)
	Token 		: 式の最初のトークン（上記の例の x）
	Expression 	: 最初のトークンに続く式（上記の例の + 10）
 */
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

/*
	Node インターフェースと Statement インターフェースを override.
 */
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/*
	整数リテラルを表す構造体型.
	Token : 整数リテラルを表すトークン
	Value : 整数リテラル
 */
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

/*
	Node インターフェースと Statement インターフェースを override.
 */
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

/*
	前置演算子を含む式の構造体型.（ex. <prefix operator><expression>;）
	Token		: 前置演算子を表すトークン（上記の <prefix operator> ex.「!」）
	Operator	: 前置演算子の文字列
	Right 		: 前置演算子の右側の式（上記の <expression>）
 */
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

/*
	Node インターフェースと Statement インターフェースを override.
 */
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
	中置演算子を含む式の構造体型.（ex. <expression> <infix operator> <expression>）
	Token		: 中置演算子を表すトークン（上記の <prefix operator> ex.「!」）
	Left		: 演算子の左側の式
	Operator	: 演算子の文字列
	Right 		: 演算子の右側の式
 */
type InfixExpression struct {
	Token    token.Token // 演算子を表すトークン（ex.「+」）
	Left     Expression
	Operator string
	Right    Expression
}

/*
	Node インターフェースと Statement インターフェースを override.
 */
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
	真偽値リテラルの構造体型.
	Token	: 真偽値リテラルを表すトークン
	Value 	: 真偽値リテラル
 */
type Boolean struct {
	Token token.Token
	Value bool
}

/*
	Node インターフェースと Expression インターフェースを override.
 */
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token // 'if' トークン
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // { トークン
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // 'fn' トークン
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // '(' トークン
	Function  Expression  // Identifier か FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
