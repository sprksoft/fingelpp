package subparsers

import (
	"fingelpp/parsermaker"
	"strings"
)

// Always put this parser at the bottom of multiParsers to prevent stealing of other parser's lines.
type ParagraphParser struct {
	builder *strings.Builder
	styler  parsermaker.InlineStyler
}

func NewParagraphParser(builder *strings.Builder) *ParagraphParser {
	return &ParagraphParser{
		builder: builder,
		styler:  BasicStyler,
	}
}
func NewParagraphParserWithStyler(builder *strings.Builder, styler parsermaker.InlineStyler) *ParagraphParser {
	return &ParagraphParser{
		builder,
		styler,
	}
}

func (*ParagraphParser) Wanted(line string) bool {
	return line != ""
}

func (p *ParagraphParser) Init() {
	p.builder.WriteString("<p>")
}

func (p *ParagraphParser) Next(line string) bool {
	if line == "" {
		return false
	}
	p.builder.WriteString(p.styler(line))
	p.builder.WriteString("<br>")
	return true
}
func (p *ParagraphParser) Finalize() {
	p.builder.WriteString("</p>")
}
