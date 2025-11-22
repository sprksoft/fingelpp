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

	p.builder.WriteString(`
	<div class="exr-footer">
		<div class="score-line">
			<span class="score"></span>
			<span class="score-bar">
				<span class="fill"></span>
			</span>
			<svg class="star" version="1.1" viewBox="0 0 75.88 72.259" xmlns="http://www.w3.org/2000/svg">
			<path d="m65.107 72.59-23.162-12.122-23.116 12.208 4.3709-25.774-18.754-18.212 25.863-3.8077 11.526-23.464 11.614 23.421 25.878 3.7107-18.686 18.283z" fill="currentColor" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
			</svg>
		</div>
		<div class="car-buttons">
			<button></button>
			<button></button>
			<button class="active"></button>
			<button></button>
		</div>
	</div>
	`)
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
