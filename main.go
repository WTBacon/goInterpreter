package main

import (
	"fmt"
	"github.com/WTBacon/goInterpreter/repl"
	"os"
	"os/user"
)

/*
	挨拶をして, インタラクティブモードスタート.
 */
func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Bacon programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
