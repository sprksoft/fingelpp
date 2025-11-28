package parsermaker

type InlineStyler func(string) string

type Parser interface {
	// Called when parser is created
	Init()

	// Check if the line wants to be parsed by the parser
	Wanted(line string) bool

	// Called for every line,
	// return false to stop parsing and give the current line to another parser
	Next(line string) bool

	// Called when the parser will be switched
	Finalize()
}
