package main

import (
	"strings"
	"fmt"
	"os"
	"io"

	"github.com/peterh/liner"
)

type HistoryMode int
const (
	HST_READ HistoryMode = iota
	HST_WRITE
)

func open_history(rl *liner.State, file_path string, mode HistoryMode) {
	f, err := os.OpenFile(file_path, os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		if mode == HST_READ {
			rl.ReadHistory(f)
		} else if mode == HST_WRITE {
			rl.WriteHistory(f)
		}
		f.Close()
	}
}

func is_null_or_whitespace(str string) bool {
	return strings.TrimSpace(str) == ""
}

/* Global access since other functions may want to prompt */
/* TODO: Name this better ... it's a global ... c'mon */
var rl = liner.NewLiner()
func main() {
	/* TODO: This should be ~/.something */
	const history_path = "history"

	defer rl.Close()
	rl.SetCtrlCAborts(true)

	open_history(rl, history_path, HST_READ)

	var parsed_line *ParsedLine
	for true {
		line, err := rl.Prompt("$ ")

		if err == liner.ErrPromptAborted {
			println("<C-c>")
			break
		} else if err == io.EOF {
			println("<C-d>")
			break
		}

		/* If an empty line is input do nothing */
		if (is_null_or_whitespace(line)) {
			continue
		}

		parsed_line = construct_parsed_line(line)
		rl.AppendHistory(line)
		fmt.Println(parsed_line.ToStringRecursive())
	}

	open_history(rl, history_path, HST_WRITE)
}
