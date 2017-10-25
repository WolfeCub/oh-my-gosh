package main

import (
	"unicode"
	"fmt"
    "strings"
)

type ConType int
const (
    SEQUENCE ConType = iota
    AND
    OR
)

func pcon(tok ConType) string {
	switch (tok) {
	case AND:
		return("AND")
	case OR:
		return("OR")
	default:
		return("SEQUENCE")
	}
}

type Token int
const (
    IDENTIFIER Token = iota
    FILEIN
    FILEOUT
    FILEOUTAPPEND
    FILEOUTDOUBLE
    /* From here on it should be a new
     * ParsedLine rather than a new PipeLine
     */
    SEMICOLON
    AMPERSAND
    PIPE
    TWOAMPERSANDS
    TWOPIPES
    DOUBLEPIPE
    EOL
)

func ptok(tok Token) string {
	switch (tok) {
	case FILEIN:
		return("FILEIN")
	case FILEOUT:
		return("FILEOUT")
	case FILEOUTAPPEND:
		return("FILEOUTAPPEND")
	case FILEOUTDOUBLE:
		return("FILEOUTDOUBLE")
	case SEMICOLON:
		return("SEMICOLON")
	case PIPE:
		return("PIPE")
	case AMPERSAND:
		return("AMPERSAND")
	case TWOAMPERSANDS:
		return("TWOAMPERSANDS")
	case TWOPIPES:
		return("TWOPIPES")
	case DOUBLEPIPE:
		return("DOUBLEPIPE")
	case EOL:
		return("EOL")
	default:
		return("ID")
	}
}

type PipeLine struct {
	argv []string
	is_double_redirect bool
    next *PipeLine
}

func (pl *PipeLine) ToString(indent bool) string {
	tab := "    "
	total := ""
	total += tab + "argv: " + fmt.Sprint(pl.argv) + "\n"
	total += tab + "is_double_redirect: " + fmt.Sprint(pl.is_double_redirect) + "\n"
	total += tab + "next: " + fmt.Sprintf("%p", pl.next)
	return total
}

func (pl *PipeLine) ToStringRecursive(indent bool) string {
	total := pl.ToString(indent)

	if pl.next != nil {
		return total + "\n\n"  + pl.next.ToStringRecursive(indent)
	} else {
		return total
	}
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

func (pl *ParsedLine) ToString() string {
	total := ""
	total += "con_type: " + pcon(pl.con_type) + "\n"
	if pl.input != nil {
		total += "input: " + *(pl.input) + "\n"
	} else {
		total += "input: nil" + "\n"
	}
	if pl.output != nil {
		total += "output: " + *(pl.output) + "\n"
	} else {
		total += "output: nil" + "\n"
	}
	total += "is_doubled: " + fmt.Sprint(pl.is_doubled) + "\n"
	total += "background: " + fmt.Sprint(pl.background) + "\n"
	total += "pipeline: \n" + pl.pipeline.ToStringRecursive(true) + "\n"
	total += "next: " + fmt.Sprintf("%p", pl.next)
	return total
}

func (pl *ParsedLine) ToStringRecursive() string {
	total := pl.ToString()

	if pl.next != nil {
		return total + "\n\n" + pl.next.ToStringRecursive() 
	} else {
		return total
	}
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
	curline.con_type = SEQUENCE
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
		for (tok < SEMICOLON) {
			switch tok {
			case IDENTIFIER:
				argv = append(argv, val)
			case FILEIN:
				if (curline.input != nil) {
					println("Error: multiple input redirects")
					return nil
				}
				tok, val = get_token(line, index)
				if (tok != IDENTIFIER) {
					println("Error: error in input redirect")
					return nil
				}
				curline.input = new(string)
				*curline.input = val
			case FILEOUTDOUBLE:
				curline.is_doubled = true
				fallthrough
			case FILEOUT:
				if (curline.output != nil) {
					println("Error: multiple input redirects")
					return nil
				}
				tok, val = get_token(line, index)
				if (tok != IDENTIFIER) {
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

		if (tok == AMPERSAND) {
			curline.background = true
		}

		if (tok == DOUBLEPIPE) {
			is_double = true
		}

		if (tok == SEMICOLON || tok == AMPERSAND || tok == TWOAMPERSANDS || tok == DOUBLEPIPE) {
			curline.next = new(ParsedLine)
			curline = curline.next

			if (tok == SEMICOLON || tok  == AMPERSAND) {
				curline.con_type = SEQUENCE
			} else if (tok == TWOAMPERSANDS) {
				curline.con_type = AND
			} else {
				curline.con_type = OR
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
		return FILEIN, content
	case '>':
		*index++
		if (line[*index] == '&') {
			*index++
			return FILEOUTDOUBLE, content
		}
		return FILEOUT, content
	case ';':
		*index++
		return SEMICOLON, content
	case '|':
		if (line[*index + 1] == '|') {
			*index += 2
			return TWOPIPES, content
		}
		*index++
		if (line[*index] == '&') {
			*index++
			return DOUBLEPIPE, content
		}
		return PIPE, content
	case '&':
		if (line[*index + 1] == '&') {
			*index += 2
			return TWOAMPERSANDS, content
		} else {
			*index++
			return AMPERSAND, content
		}
	}
	/* If we've reached here we know it's an identifier */
	start := *index
	for (*index < len(line) && !unicode.IsSpace(rune(line[*index])) && strings.IndexByte("<>;&|", line[*index]) == -1) {
		*index++
	}
	
	return IDENTIFIER, line[start:*index]
}
