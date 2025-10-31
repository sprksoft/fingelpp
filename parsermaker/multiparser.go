package parsermaker

import "strings"

func NewMultiParser(builder *strings.Builder, parsers []Parser) *MultiParser {
	return &MultiParser{
		curParser: nil,
		builder:   builder,
		parsers:   parsers,
	}
}

type MultiParser struct {
	curParser Parser
	builder   *strings.Builder
	parsers   []Parser
}

func (p *MultiParser) switchParser(new Parser) {
	if p.curParser != nil {
		p.curParser.Finalize()
	}
	p.curParser = new

	if p.curParser != nil {
		p.curParser.Init()
	}
}

func (p *MultiParser) chooseParser(line string) Parser {
	for _, par := range p.parsers {
		if par.Wanted(line) {
			return par
		}
	}
	return nil
}

func (p *MultiParser) Wanted(line string) bool {
	for _, par := range p.parsers {
		if par.Wanted(line) {
			return true
		}
	}
	return false
}

func (p *MultiParser) Init() {}

func (p *MultiParser) Next(line string) bool {
	if p.curParser == nil || !p.curParser.Next(line) {
		par := p.chooseParser(line)

		p.switchParser(par)

		if p.curParser == nil && line != "" { // Allow empty lines to be unparsed
			// Can't fin a parser to parse the current line.
			return false
		}
		if p.curParser != nil && !p.curParser.Next(line) {
			panic("BUG: parsed didn't want to parse line that it said it wanted in the wanted function.")
		}
	}
	return true
}

func (p *MultiParser) Finalize() {
	if p.curParser != nil {
		p.curParser.Finalize()
	}
}
