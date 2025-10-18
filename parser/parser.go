package parser

import (
	"html/template"
	"regexp"
	"strconv"
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
	init(builder *strings.Builder)

	// Called for every line,
	// return false to stop parsing and give the current line to another parser
	next(builder *strings.Builder, line string) bool

	// Called when the parser will be switched
	finalize(builder *strings.Builder)
}

type paragraphParser struct{}

func (p paragraphParser) init(builder *strings.Builder) {
	builder.WriteString("<p>")
}

func (p paragraphParser) next(builder *strings.Builder, line string) bool {
	if line == "" {
		return false
	}
	builder.WriteString(finalizeTextBlock(line))
	builder.WriteString("<br>")
	return true
}
func (p paragraphParser) finalize(builder *strings.Builder) {
	builder.WriteString("</p>")
}

type listParser struct{}

func (p listParser) init(builder *strings.Builder) {
	builder.WriteString("<ul>")
}

func (p listParser) next(builder *strings.Builder, line string) bool {
	if !strings.HasPrefix(line, "-") {
		return false
	}
	content := strings.TrimSpace(line[1:])
	builder.WriteString("<li>")
	builder.WriteString(finalizeTextBlock(content))
	builder.WriteString("</li>")
	return true
}
func (p listParser) finalize(builder *strings.Builder) {
	builder.WriteString("</ul>")
}

type titleParser struct{}

func (p titleParser) init(builder *strings.Builder) {}

func (p titleParser) next(builder *strings.Builder, line string) bool {
	for i, char := range line {
		if char != '#' {
			if i == 0 { // no hashtag found
				return false
			}
			title := strings.TrimSpace(line[i:])
			headingNum := strconv.FormatInt(int64(i+1), 10)
			builder.WriteString("<h" + headingNum + ">" + title + "</h" + headingNum + ">")
			return true
		}
	}
	return false
}
func (p titleParser) finalize(builder *strings.Builder) {}

func ParseFinSyn(mdText string) template.HTML {
	lines := strings.Split(mdText, "\n")
	var builder strings.Builder

	var curParser parser
	curParser = paragraphParser{}

	switchParser := func(new parser) {
		curParser.finalize(&builder)
		curParser = new
		curParser.init(&builder)
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "- ") {
			switchParser(listParser{})
		}

		if strings.HasPrefix(line, "#") {
			switchParser(titleParser{})
		}

		if !curParser.next(&builder, line) {
			// When the current parser doesn't want to parse. we switch to the paragraph parser
			switchParser(paragraphParser{})
			curParser.next(&builder, line)
		}
	}

	curParser.finalize(&builder)

	return template.HTML(builder.String())
}
