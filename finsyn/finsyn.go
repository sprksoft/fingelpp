package finsyn

import (
	"bufio"
	"fingelpp/finsyn/subparsers"
	"github.com/charmbracelet/log"
	"html/template"
	"strings"
)

func ParseFinSyn(mdText string) template.HTML {
	var builder strings.Builder

	parser := subparsers.NewFinSynParser(&builder)
	parser.Init()

	scanner := bufio.NewScanner(strings.NewReader(mdText))
	i := 0
	for scanner.Scan() {
		line, _, _ := strings.Cut(scanner.Text(), "// ") // strip comments

		if !parser.Next(line) {
			log.Errorf("Failed to parse insyn on line %v: %v\n", i+1, line)
			break
		}
		i++
	}

	parser.Finalize()

	return template.HTML(builder.String())
}
