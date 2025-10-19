package parser

import (
	"strconv"
	"strings"
)

type paragraphParser struct{ builder *strings.Builder }

func (p *paragraphParser) init() {
	p.builder.WriteString("<p>")
}

func (p *paragraphParser) next(line string) bool {
	if line == "" {
		return false
	}
	p.builder.WriteString(finalizeTextBlock(line))
	p.builder.WriteString("<br>")
	return true
}
func (p *paragraphParser) finalize() {
	p.builder.WriteString("</p>")
}

type listParser struct{ builder *strings.Builder }

func (p *listParser) init() {
	p.builder.WriteString("<ul>")
}

func (p *listParser) next(line string) bool {
	if !strings.HasPrefix(line, "-") {
		return false
	}
	content := strings.TrimSpace(line[1:])
	p.builder.WriteString("<li>")
	p.builder.WriteString(finalizeTextBlock(content))
	p.builder.WriteString("</li>")
	return true
}
func (p *listParser) finalize() {
	p.builder.WriteString("</ul>")
}

type titleParser struct{ builder *strings.Builder }

func (p *titleParser) init() {}

func (p *titleParser) next(line string) bool {
	for i, char := range line {
		if char != '#' {
			if i == 0 { // no hashtag found
				return false
			}
			title := strings.TrimSpace(line[i:])
			headingNum := strconv.FormatInt(int64(i+1), 10)
			p.builder.WriteString("<h" + headingNum + ">" + title + "</h" + headingNum + ">")
			return true
		}
	}
	return false
}
func (p *titleParser) finalize() {}

type infoParser struct {
	builder *strings.Builder
	finsyn  finSynParser
}

func (p *infoParser) init() {
	p.builder.WriteString("<section class=\"block info\">")
	p.finsyn = finSynParser{builder: p.builder}
	p.finsyn.init()
}

func (p *infoParser) next(line string) bool {
	if !strings.HasPrefix(line, ">") {
		return false
	}
	line = strings.TrimSpace(line[1:])

	if strings.HasPrefix(line, "[INFO]") {
		title := strings.TrimSpace(line[len("[INFO]"):])
		p.builder.WriteString("<h1 class=\"block-title\">")
		p.builder.WriteString(title)
		p.builder.WriteString("</h1>")
	} else {
		p.finsyn.next(line)
	}

	return true
}

func (p *infoParser) finalize() {
	p.finsyn.finalize()
	p.builder.WriteString("</section>")
}

type exersiseParser struct {
	builder *strings.Builder
	finsyn  finSynParser
}

func (p *exersiseParser) init() {
	p.builder.WriteString("<section class=\"block exercise\">")
	p.finsyn = finSynParser{builder: p.builder}
	p.finsyn.init()
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
	} else {
		p.finsyn.next(line)
	}

	return true
}

func (p *exersiseParser) finalize() {
	p.finsyn.finalize()
	p.builder.WriteString("</section>")
}
