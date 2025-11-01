package subparsers

import (
	"fingelpp/parsermaker"
	"strings"
)

type InfoParser struct {
	builder *strings.Builder
	parser  parsermaker.MultiParser
}

func NewInfoParser(builder *strings.Builder) *InfoParser {
	parser := parsermaker.NewMultiParser(builder, []parsermaker.Parser{
		NewListParser(builder),
		NewParagraphParser(builder),
	})

	return &InfoParser{
		builder: builder,
		parser:  *parser,
	}
}

func (*InfoParser) Wanted(line string) bool {
	return strings.HasPrefix(line, "> [INFO]")
}

func (p *InfoParser) Init() {
	p.builder.WriteString("<section class=\"block info\">")
	p.parser.Init()
}

func (p *InfoParser) Next(line string) bool {
	if !strings.HasPrefix(line, ">") {
		return false
	}
	line = strings.TrimSpace(line[1:])

	if strings.HasPrefix(line, "[INFO]") {
		title := strings.TrimSpace(line[len("[INFO]"):])
		p.builder.WriteString("<div class=block-title><h1>")
		p.builder.WriteString(title)
		p.builder.WriteString("</h1></div>")
		return true
	} else {
		return p.parser.Next(line)
	}
}

func (p *InfoParser) Finalize() {
	p.parser.Finalize()
	p.builder.WriteString("</section>")
}
