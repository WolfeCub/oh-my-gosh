package main

import (
	"strings"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"io/ioutil"

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

func shell_completer(line string) (c []string) {
	split := strings.Split(line, " ")
	word_at_point := split[len(split)-1]

	dir, base := filepath.Split(word_at_point)

	files, err := ioutil.ReadDir(dir)
	if (err == nil) {
		for _, file := range files {
			if strings.HasPrefix(file.Name(), base) {
				/* TODO: Figure out if there's a more efficient way to do this */
				c = append(c, dir + file.Name())
			}
		}
	}
	return
}

/* Global access since other functions may want to prompt */
/* TODO: Name this better ... it's a global ... c'mon */
var rl = liner.NewLiner()
func main() {
	/* TODO: This should be ~/.something */
	const history_path = "history"

	defer rl.Close()
	rl.SetCtrlCAborts(true)
	rl.SetCompleter(shell_completer)
	rl.SetTabCompletionStyle(liner.TabPrints)

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
