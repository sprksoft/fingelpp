package subparsers

import (
	"fingelpp/parsermaker"
	"strings"
)

type ListParser struct {
	builder *strings.Builder
	styler  parsermaker.InlineStyler
}

func NewListParser(builder *strings.Builder) *ListParser {
	return &ListParser{
		builder: builder,
		styler:  BasicStyler,
	}
}
func NewListParserWithStyler(builder *strings.Builder, styler parsermaker.InlineStyler) *ListParser {
	return &ListParser{
		builder,
		styler,
	}
}

func (*ListParser) Wanted(line string) bool {
	return strings.HasPrefix(line, "- ")
}

func (p *ListParser) Init() {
	p.builder.WriteString("<ul>")
}

func (p *ListParser) Next(line string) bool {
	if !strings.HasPrefix(line, "-") {
		return false
	}
	content := strings.TrimSpace(line[1:])
	p.builder.WriteString("<li>")
	p.builder.WriteString(p.styler(content))
	p.builder.WriteString("</li>")
	return true
}
func (p *ListParser) Finalize() {
	p.builder.WriteString("</ul>")
}
