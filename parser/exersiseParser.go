package parser

import (
	"strings"
)

type exersiseParser struct {
	builder   *strings.Builder
	curParser finSynParser
}

func (*exersiseParser) wanted(line string) bool {
	return strings.HasPrefix(line, "> [EX]")
}

func (p *exersiseParser) init() {
	p.builder.WriteString("<section class=\"block exercise\">")
	p.curParser = finSynParser{builder: p.builder}
	p.curParser.init()
}

func (p *exersiseParser) next(line string) bool {
	if !strings.HasPrefix(line, ">") {
		return false
	}
	line = strings.TrimSpace(line[1:])

	if strings.HasPrefix(line, "[EX]") {
		title := strings.TrimSpace(line[len("[INFO]"):])
		p.builder.WriteString("<h1 class=\"block-title\">")
		p.builder.WriteString(title)
		p.builder.WriteString("</h1>")
	} else if (*multipleChoiceParser).wanted(nil, line) {
		p.curParser.finalize()
	}

	p.curParser.next(line)

	return true
}

func (p *exersiseParser) finalize() {
	p.curParser.finalize()
	p.builder.WriteString("</section>")
}

type multipleChoiceParser struct {
	builder *strings.Builder
}

func (*multipleChoiceParser) wanted(line string) bool {
	return strings.HasPrefix(line, "@[x]") || strings.HasPrefix(line, "@[o]")
}

func (p *multipleChoiceParser) init() {
	p.builder.WriteString("</ul class=\"exr-multiplechoice>")
}

func (p *multipleChoiceParser) next(line string) bool {

	var awnser bool
	if strings.HasPrefix(line, "@[x]") {
		awnser = true
	} else if strings.HasPrefix(line, "@[o]") {
		awnser = true
	} else {
		return false
	}

	p.builder.WriteString("<li><input type=checkbox data-awnser=")
	if awnser {
		p.builder.WriteString("true")
	} else {
		p.builder.WriteString("false")
	}
	p.builder.WriteString("><label>")
	p.builder.WriteString(line[len("@[x]"):])
	p.builder.WriteString("</label></li>")

	return true
}

func (p *multipleChoiceParser) finalize() {
	p.builder.WriteString("</ul>")
}
