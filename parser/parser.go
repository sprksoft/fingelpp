package parser

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// Applies the inline styles (bold, italic, underline,...)
func finalizeTextBlock(text string) string {
	boldItalicRe := regexp.MustCompile(`\*\*\_(.*?)\_\*\*`) // **_text_**
	boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)           // **text**
	italicRe := regexp.MustCompile(`\_(.*?)\_`)             // _text_

	text = boldItalicRe.ReplaceAllString(text, "<strong><em>$1</em></strong>")
	text = boldRe.ReplaceAllString(text, "<strong>$1</strong>")
	text = italicRe.ReplaceAllString(text, "<em>$1</em>")
	return text
}

type parser interface {
	// Called when parser is created
	init()

	// Called for every line,
	// return false to stop parsing and give the current line to another parser
	next(line string) bool

	// Called when the parser will be switched
	finalize()

	// Check if the line wants to be parsed by the parser
	wanted(line string) bool
}

type finSynParser struct {
	curParser parser
	builder   *strings.Builder
}

func (p *finSynParser) switchParser(new parser) {
	if p.curParser != nil {
		p.curParser.finalize()
	}
	p.curParser = new
	p.curParser.init()
}

func (p *finSynParser) init() {}

func (p *finSynParser) next(line string) bool {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "//") { // skip comments
		return true
	}

	if (*listParser).wanted(nil, line) {
		p.switchParser(&listParser{builder: p.builder})
	}

	if (*titleParser).wanted(nil, line) {
		p.switchParser(&titleParser{builder: p.builder})
	}

	if (*infoParser).wanted(nil, line) {
		p.switchParser(&infoParser{builder: p.builder})
	}
	if (*exersiseParser).wanted(nil, line) {
		p.switchParser(&exersiseParser{builder: p.builder})
	}

	if (p.curParser == nil && line != "") || !p.curParser.next(line) {
		// Switch to the paragraph parser. when current parser doesn't want to parse or when parser is nil
		p.switchParser(&paragraphParser{builder: p.builder})
		p.curParser.next(line)
	}

	return true
}
func (p *finSynParser) finalize() {
	p.curParser.finalize()
}

func ParseFinSyn(mdText string) template.HTML {
	lines := strings.Split(mdText, "\n")
	var builder strings.Builder

	parser := finSynParser{builder: &builder}
	parser.init()
	fmt.Printf("%v\n", parser.curParser)

	for _, line := range lines {
		parser.next(line)
	}

	parser.finalize()

	return template.HTML(builder.String())
}
