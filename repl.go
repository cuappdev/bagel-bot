package main

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/kong"
	"gorm.io/gorm"
	"os"
	"strings"
)

func printRepl() {
	fmt.Print("bagel > ")
}

func get(r *bufio.Reader) string {
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}

func shouldContinue(text string) bool {
	if strings.EqualFold("exit", text) {
		return false
	}
	return true
}

func BagelRepl(db *gorm.DB, s *Slack) {
	reader := bufio.NewReader(os.Stdin)
	printRepl()
	text := get(reader)
	for ; shouldContinue(text); text = get(reader) {
		if err := Run(text, os.Stdout, os.Stderr, db, s); err != nil {
			if parseError, success := err.(kong.ParseError); success {
				fmt.Println(parseError)
			} else {
				panic(err)
			}
		}
		printRepl()
	}
}