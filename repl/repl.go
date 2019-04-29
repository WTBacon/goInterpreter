package repl

import (
	"bufio"
	"fmt"
	"github.com/WTBacon/goInterpreter/lexer"
	"github.com/WTBacon/goInterpreter/token"
	"io"
)

const PROMPT = ">> "

/*
	入力の読み込み（Read）, インタプリタに送って評価（Eval）,
	インタプリタの結果/出力を表示（print）, を繰り返す（Loop）.
 */
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
