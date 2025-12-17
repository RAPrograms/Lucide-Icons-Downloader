package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type icon_details struct {
	Name         string
	Contributors []string
	SVG          string
}

func (details icon_details) format_svg() (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(details.SVG))
	decoder.Strict = false

	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}
		if err := encoder.EncodeToken(tok); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (details icon_details) format_contributors() string {
	output := details.Contributors[0]
	for i := 1; i < len(details.Contributors); i++ {
		if i >= len(details.Contributors)-1 {
			output += " & "
		} else {
			output += ", "
		}

		output += details.Contributors[i]
	}
	return output
}

func (details icon_details) to_string() string {
	svg, svg_err := details.format_svg()
	if svg_err != nil {
		// TODO: add flag to force format
		svg = details.SVG
	}

	return fmt.Sprintf("<!-- Icon sourced from Lucide.dev -->\n"+
		"<!-- Name: %s -->\n"+
		"<!-- Contributors: %s -->\n\n%s",
		details.Name, details.format_contributors(), svg,
	)
}
