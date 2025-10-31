package subparsers

import (
	"strconv"
	"strings"
)

func NewTitleParser(builder *strings.Builder) *TitleParser {
	return &TitleParser{builder}
}

type TitleParser struct{ builder *strings.Builder }

func (*TitleParser) Wanted(line string) bool {
	return strings.HasPrefix(line, "#")
}

func (p *TitleParser) Init() {}

func (p *TitleParser) Next(line string) bool {
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
func (p *TitleParser) Finalize() {}
