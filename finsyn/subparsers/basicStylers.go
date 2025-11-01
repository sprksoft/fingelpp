package subparsers

import (
	"regexp"
)

var boldRe = regexp.MustCompile(`\*\*(.*?)\*\*`)            // **text**
var italicRe = regexp.MustCompile(`\_(.*?)\_`)              // _text_
var linkRe = regexp.MustCompile(`\[([^\n]*)\]\(([^\s]*)\)`) // [text](https://google.com)

func BoldStyler(text string) string {
	return boldRe.ReplaceAllString(text, "<strong>$1</strong>")
}

func ItalicStyler(text string) string {
	return italicRe.ReplaceAllString(text, "<em>$1</em>")
}

func linksStyler(text string) string {
	return linkRe.ReplaceAllString(text, "<a href=\"$2\">$1</a>")
}

func BasicStyler(s string) string {
	s = BoldStyler(s)
	s = ItalicStyler(s)
	s = linksStyler(s)
	return s
}
