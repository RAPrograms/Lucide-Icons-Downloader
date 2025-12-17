package main

import "fmt"

type icon_details struct {
	Name         string
	Contributors []string
	SVG          string
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
	return fmt.Sprintf("<!-- Icon sourced from Lucide.dev -->\n"+
		"<!-- Name: %s -->\n"+
		"<!-- Contributors: %s -->\n\n%s",
		details.Name, details.format_contributors(), details.SVG,
	)
}
