package subparsers

import (
	"fingelpp/parsermaker"
	"strings"
)

type FinSynParser struct {
	builder *strings.Builder
	parser  parsermaker.MultiParser
}

func NewFinSynParser(builder *strings.Builder) *FinSynParser {
	parser := parsermaker.NewMultiParser(builder, []parsermaker.Parser{
		NewTitleParser(builder),
		NewListParser(builder),
		NewInfoParser(builder),
		NewExersizeParser(builder),
		NewParagraphParser(builder),
	})

	return &FinSynParser{
		builder: builder,
		parser:  *parser,
	}
}

func (p *FinSynParser) Wanted(line string) bool {
	return p.parser.Wanted(line)
}

func (p *FinSynParser) Init() {
	p.parser.Init()
}

func (p *FinSynParser) Next(line string) bool {
	return p.parser.Next(line)
}
func (p *FinSynParser) Finalize() {
	p.parser.Finalize()
}
