package main

import (
	"unicode"
	"fmt"
	"bufio"
	"strings"
	"os"
)

type ConType int
const (
    Sequence ConType = iota
    And
    Or
)

type Token int
const (
    Identifier Token = iota
    FileIn
    FileOut
    FileOutAppend
    FileOutDouble
    /* From here on it should be a new
     * ParsedLine rather than a new PipeLine
     */
    Semicolon
    Ampersand
    Pipe
    TwoAmpersands
    TwoPipes
    DoublePipe
    EOL
)

type PipeLine struct {
	argv []string
	is_double_redirect bool
    next *PipeLine
}

type ParsedLine struct {
	con_type ConType
	input *string
	output *string
	is_doubled bool
	background bool
	pipeline *PipeLine
	next * ParsedLine
}

func construct_parsed_line (line string) *ParsedLine {
	var retval *ParsedLine
	var curline *ParsedLine
	var plp **PipeLine
	var argv []string
	is_double := false
	var tok Token
	var val string

	retval = new(ParsedLine)
	curline = retval
	curline.con_type = Sequence
	curline.input = nil
	curline.output = nil
	curline.is_doubled = false
	curline.pipeline = nil
	curline.next = nil
	plp = &(curline.pipeline)

	var index *int = new(int)
	*index = 0
	
	tok, val = get_token(line, index)
	for (tok != EOL) {
		for (tok < Semicolon) {
			switch tok {
			case Identifier:
				argv = append(argv, val)
			case FileIn:
				if (curline.input != nil) {
					println("Error: multiple input redirects")
					return nil
				}
				tok, val = get_token(line, index)
				if (tok != Identifier) {
					println("Error: error in input redirect")
					return nil
				}
				curline.input = new(string)
				*curline.input = val
			case FileOutDouble:
				curline.is_doubled = true
				fallthrough
			case FileOut:
				if (curline.output != nil) {
					println("Error: multiple input redirects")
					return nil
				}
				tok, val = get_token(line, index)
				if (tok != Identifier) {
					println("Error: error in output redirect")
					return nil
				}
				curline.output = new(string)
				*curline.output = val
			}
			tok, val = get_token(line, index)
		}


		if (len(argv) > 0) {
			*plp = new(PipeLine)
			(*plp).next = nil
			(*plp).argv = argv
			(*plp).is_double_redirect = is_double
			plp = &((*plp).next)
			is_double = false
			argv = nil
		} else if (tok != EOL) {
			println("Error: null command found before " + ptok(tok))
			return nil
		}

		if (tok == Ampersand) {
			curline.background = true
		}

		if (tok == DoublePipe) {
			is_double = true
		}

		if (tok == Semicolon || tok == Ampersand || tok == TwoAmpersands || tok == DoublePipe) {
			curline.next = new(ParsedLine)
			curline = curline.next

			if (tok == Semicolon || tok  == Ampersand) {
				curline.con_type = Sequence
			} else if (tok == TwoAmpersands) {
				curline.con_type = And
			} else {
				curline.con_type = Or
			}

			curline.input = nil
			curline.output = nil
			curline.is_doubled = false
			curline.background = false
			curline.pipeline = nil
			curline.next = nil
			plp = &(curline.pipeline)
		}

		tok, val = get_token(line, index)
	}

	return retval
}

func get_token(line string, index *int) (Token, string){
	var content string

	for (*index < len(line) && unicode.IsSpace(rune(line[*index]))) {
		*index++;
	}

	if (*index >= len(line)) {
		return EOL, content
	}

	switch line[*index] {
	case '<':
		*index++
		return FileIn, content
	case '>':
		*index++
		if (line[*index] == '&') {
			*index++
			return FileOutDouble, content
		}
		return FileOut, content
	case ';':
		*index++
		return Semicolon, content
	case '|':
		if (line[*index + 1] == '|') {
			*index += 2
			return TwoPipes, content
		}
		*index++
		if (line[*index] == '&') {
			*index++
			return DoublePipe, content
		}
		return Pipe, content
	case '&':
		if (line[*index + 1] == '&') {
			*index += 2
			return TwoAmpersands, content
		} else {
			*index++
			return Ampersand, content
		}
	}
	/* If we've reached here we know it's an identifier */
	start := *index
	for (*index < len(line) && !unicode.IsSpace(rune(line[*index])) && strings.IndexByte("<>;&|", line[*index]) == -1) {
		*index++
	}
	
	return Identifier, line[start:*index]
}

func ptok(tok Token) string {
	switch (tok) {
	case FileIn:
		return("FileIn");
	case FileOut:
		return("FileOut");
	case Semicolon:
		return("Semicolon");
	case Pipe:
		return("Pipe");
	case Ampersand:
		return("Ampersand");
	case TwoAmpersands:
		return("TwoAmpersands");
	case TwoPipes:
		return("TwoPipes");
	case FileOutDouble:
		return("FileOutDouble");
	case DoublePipe:
		return("DoublePipe");
	case EOL:
		return("EOL");
	default:
		return("ID");
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var parsed_line *ParsedLine
	for true {
		fmt.Print("$ ")
		text, _ := reader.ReadString('\n')
		if (len(strings.TrimSpace(text)) == 0) {continue}
		parsed_line = construct_parsed_line(text)
		p := parsed_line.pipeline
		for (p != nil) {
			fmt.Println(p.argv)
			p = p.next
		}
	}
}
