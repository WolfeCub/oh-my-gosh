package main

import "os"

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

type Pipeline struct {
	argv []string
	is_double_redirect bool
    next *Pipeline
}

type ParsedLine struct {
	con_type Token
	input *os.File
	output *os.File
	is_doubled bool
	background bool
	pipeline *Pipeline
	next * ParsedLine
}

func main() {
	println("foob")
}
