package finsyn

import (
	"fingelpp/finsyn/subparsers"
	"github.com/charmbracelet/log"
	"html/template"
	"strings"
)

func ParseFinSyn(mdText string) template.HTML {
	var builder strings.Builder

	parser := subparsers.NewFinSynParser(&builder)
	parser.Init()

	lines := strings.Split(mdText, "\n")
	for i, line := range lines {
		line, _, _ = strings.Cut(line, "// ") // strip comments
		line = strings.TrimSpace(line)

		if !parser.Next(line) {
			log.Errorf("Failed to parse insyn on line %v: %v\n", i+1, line)
			break
		}
	}

	parser.Finalize()

	return template.HTML(builder.String())
}
