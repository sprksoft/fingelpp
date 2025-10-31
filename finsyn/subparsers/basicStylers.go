package subparsers

import (
	"regexp"
)

var boldRe = regexp.MustCompile(`\*\*(.*?)\*\*`) // **text**
var italicRe = regexp.MustCompile(`\_(.*?)\_`)   // _text_

func BoldStyler(text string) string {
	return boldRe.ReplaceAllString(text, "<strong>$1</strong>")
}

func ItalicStyler(text string) string {
	return italicRe.ReplaceAllString(text, "<em>$1</em>")
}

func BasicStyler(s string) string {
	s = BoldStyler(s)
	s = ItalicStyler(s)
	return s
}
