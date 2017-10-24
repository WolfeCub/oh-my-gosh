package main

import (
	"strings"
	"fmt"
	"os"

	"github.com/peterh/liner"
)

func main() {
	rl := liner.NewLiner()
	defer rl.Close()
	rl.SetCtrlCAborts(true)

	/* TODO: This should be ~/.something */
	if f, err := os.Open("./history"); err == nil {
		rl.ReadHistory(f)
		f.Close()
	}

	var parsed_line *ParsedLine
	for true {
		line, err := rl.Prompt("$ ")
		if (err != nil) { panic(err) }
		if (len(strings.TrimSpace(line)) == 0) { continue }

		parsed_line = construct_parsed_line(line)
		rl.AppendHistory(line)
		fmt.Println(parsed_line.ToStringRecursive())
	}
}
