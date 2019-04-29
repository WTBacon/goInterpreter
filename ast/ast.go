package ast

import (
	"github.com/WTBacon/goInterpreter/token"
)

/*
	抽象構文木（Abstract sytax tree; AST）: ソースコードを入力とした構文解析器の出力.
	以下は, 再帰下降構文解析器のという AST の実装.
*/

/*
	AST の全てのノードは Node interface を実装しなければならない.
	つまり, TokenLiteral() メソッドを override する.
	TokenLiteral() : ノードに関連づけられているトークンのリテラル値を返す.
*/
type Node interface {
	TokenLiteral() string
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
	let 文を表す構造体型.
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
	Statement インターフェースと Node インターフェースを実装.
 */
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

/*
	識別子を表す構造体型.
	Toke 	: 識別子を示すトークン
	Value	: 識別子の値
 */
type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

/*
	Expression インターフェースと Node インターフェースを実装.
 */
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
