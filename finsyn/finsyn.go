package finsyn

import (
	"fingelpp/finsyn/subparsers"
	"github.com/charmbracelet/log"
	"html/template"
	"strings"
)

func ParseFinSyn(mdText string) template.HTML {
	lines := strings.Split(mdText, "\n")
	var builder strings.Builder

	parser := subparsers.NewFinSynParser(&builder)
	parser.Init()

	for i, line := range lines {
		line, _, _ = strings.Cut(line, "// ") // strip comments

		if !parser.Next(line) {
			log.Errorf("Failed to parse insyn on line %v: %v\n", i+1, line)
			break
		}
	}

	parser.Finalize()

	return template.HTML(builder.String())
}
