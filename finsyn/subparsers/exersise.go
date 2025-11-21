package subparsers

import (
	"fingelpp/parsermaker"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

var exerRe = regexp.MustCompile(`@(\[[a-z]+\])?\(([^()\n]*)\)`) // @[text](value)

func exerStyler(s string) string {
	s = exerRe.ReplaceAllStringFunc(s, func(s string) string {
		inputType := "text"
		s = strings.TrimSpace(s)

		var found bool
		s, found = strings.CutPrefix(s, "@[")
		if found {
			inputType, s, _ = strings.Cut(s, "](")
		} else {
			s = s[2:]
		}
		awnser := strings.TrimSuffix(s, ")")

		return "<input autocomplete=no class=\"exr\" data-type=" + inputType + " type=\"text\" data-awnser=\"" + awnser + "\">"
	})
	s = BasicStyler(s)
	return s
}

type ExersiseParser struct {
	builder *strings.Builder
	parser  parsermaker.MultiParser
}

func NewExersizeParser(builder *strings.Builder) *ExersiseParser {
	parser := parsermaker.NewMultiParser(builder, []parsermaker.Parser{
		&multipleChoiceParser{builder: builder, styler: BasicStyler, checkboxStyle: CheckboxStyleUnknown},
		NewListParser(builder),
		NewParagraphParserWithStyler(builder, exerStyler),
	})

	return &ExersiseParser{
		builder: builder,
		parser:  *parser,
	}
}

func (*ExersiseParser) Wanted(line string) bool {
	return strings.HasPrefix(line, "> [EX]")
}

func (p *ExersiseParser) Init() {
	p.builder.WriteString("<section class=\"block exercise\">")
	p.parser.Init()
}

func (p *ExersiseParser) Next(line string) bool {
	if !strings.HasPrefix(line, ">") {
		return false
	}
	line = strings.TrimSpace(line[1:])

	if strings.HasPrefix(line, "[EX]") {
		title := strings.TrimSpace(line[len("[EX]"):])
		p.builder.WriteString("<div class=\"block-title\">")
		p.builder.WriteString("<h1>")
		p.builder.WriteString(title)
		p.builder.WriteString("</h1><span class=score></span></div>")
		return true
	} else {
		return p.parser.Next(line)
	}
}

func (p *ExersiseParser) Finalize() {
	p.parser.Finalize()
	p.builder.WriteString("</section>")
}

type CheckboxStyle int

const (
	CheckboxStyleUnknown  CheckboxStyle = 0
	CheckboxStyleMultiple CheckboxStyle = 1
	CheckboxStyleSingle   CheckboxStyle = 2
)

type multipleChoiceParser struct {
	builder       *strings.Builder
	styler        parsermaker.InlineStyler
	checkboxStyle CheckboxStyle
	name          string
}

func (*multipleChoiceParser) Wanted(line string) bool {
	return strings.HasPrefix(line, "@[x]") || strings.HasPrefix(line, "@[o]") || strings.HasPrefix(line, "@(x)") || strings.HasPrefix(line, "@(o)")
}

func (p *multipleChoiceParser) Init() {
	p.builder.WriteString("<div class=\"exr-multiplechoice\"><ul>")
	p.name = strconv.Itoa(rand.Int())
}

func (p *multipleChoiceParser) Next(line string) bool {
	if !p.Wanted(line) {
		return false
	}
	if strings.HasPrefix(line, "@(") {
		p.checkboxStyle = CheckboxStyleSingle
	} else if strings.HasPrefix(line, "@[") {
		p.checkboxStyle = CheckboxStyleMultiple
	}
	awnser := line[2] == 'o'

	p.builder.WriteString("<li><label><input type=")
	switch p.checkboxStyle {
	case CheckboxStyleSingle:
		p.builder.WriteString("radio")
	case CheckboxStyleMultiple:
		p.builder.WriteString("checkbox")
	}
	p.builder.WriteString(" name=")
	p.builder.WriteString(p.name)
	p.builder.WriteString(" data-awnser=")
	if awnser {
		p.builder.WriteString("true")
	} else {
		p.builder.WriteString("false")
	}
	p.builder.WriteString("><span>")
	p.builder.WriteString(p.styler(line[len("@[x]"):]))
	p.builder.WriteString("</span></label></li>")

	return true
}

func (p *multipleChoiceParser) Finalize() {
	p.builder.WriteString("</ul>")
	p.builder.WriteString("<button>check</button>")
	p.builder.WriteString("</div>")
}
