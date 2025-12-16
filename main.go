package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

const BASE_URL = "https://lucide.dev/icons/"

func fetch_page(icon_name string) (*goquery.Document, error) {
	url := BASE_URL + strings.ToLower(icon_name)

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Network fail: " + err.Error())
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.New("Error parsing HTML: " + err.Error())
	}

	return doc, nil
}

func extract_details(page *goquery.Document) icon_details {
	// Get icon name
	name := page.Find("main h1").Text()

	// Get contributors' usernames
	contributors_selection := page.Find("main .contributors a")
	contributors := make([]string, contributors_selection.Length())

	contributors_selection.Each(func(i int, link *goquery.Selection) {
		username, _ := link.Attr("data-name")
		contributors[i] = username
	})

	// Get icon svg code
	icon_target := page.Find("main svg").First()
	svg, _ := goquery.OuterHtml(icon_target)

	return icon_details{
		Name:         name,
		Contributors: contributors,
		SVG:          svg,
	}
}

func save_icon(details icon_details) {
	// Create a new file named "myfile.txt" in the same directory
	file, err := os.Create("myfile.txt")
	if err != nil { // Check for an error during file creation
		panic(err)
	}
	defer file.Close() // Ensure the file is closed when the function exits

	_, write_err := file.WriteString(details.to_string())
	if write_err != nil {
		panic(write_err)
	}
}

func main() {
	doc, err := fetch_page("move-right")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	details := extract_details(doc)

	save_icon(details)
}
