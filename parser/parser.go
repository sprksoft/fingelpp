package parser

import (
	"regexp"
	"strings"
)

type block interface {
	getHTML()
}

type text struct {
	Content string
}

type kennis struct {
	Name, Content string
}

type exercise struct {
	Name, Content string
}

func (t text) getHTML() string {
	return t.Content
}

func (k kennis) getHTML() string {
	return k.Content
}

func (e exercise) getHTML() string {
	return e.Content
}

func getBlockHTML(b block) string {
	return b.getHTML()
}

func ParseMd(mdText string) string {
	lines := strings.Split(mdText, "\n")
	var builder strings.Builder
	insideList := false

	boldItalicRe := regexp.MustCompile(`\*\*\_(.*?)\_\*\*`) // **_text_**
	boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)           // **text**
	italicRe := regexp.MustCompile(`\_(.*?)\_`)             // _text_

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "##") {
			insideList = false

			// Extract and clean title
			content := strings.TrimPrefix(line, "##")

			builder.WriteString("<h2>" + content + "</h2>\n")

		} else if strings.HasPrefix(line, "# ") {
			insideList = false
			content := strings.TrimPrefix(line, "# ")
			builder.WriteString("<h1>" + content + "</h1>\n")

		} else if strings.HasPrefix(line, "- ") {
			line := strings.TrimPrefix(line, "- ")
			if !insideList {
				builder.WriteString("<ul>")
			}
			line = boldItalicRe.ReplaceAllString(line, "<strong><em>$1</em></strong>")
			line = boldRe.ReplaceAllString(line, "<strong>$1</strong>")
			line = italicRe.ReplaceAllString(line, "<em>$1</em>")
			builder.WriteString(`<li>` + line + `</li>`)
			insideList = true
		} else if line != "" {
			insideList = false
			line = boldItalicRe.ReplaceAllString(line, "<strong><em>$1</em></strong>")
			line = boldRe.ReplaceAllString(line, "<strong>$1</strong>")
			line = italicRe.ReplaceAllString(line, "<em>$1</em>")
			builder.WriteString("<p>" + line + "</p>\n")
		}
		if !insideList {
			builder.WriteString("</ul>")
		}
	}

	// Final list cleanup
	if insideList {
		builder.WriteString("</ul>")
	}
	return builder.String()
}
