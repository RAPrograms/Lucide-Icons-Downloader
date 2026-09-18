package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func extract_details(page *goquery.Document) (icon_details, error) {
	// Get icon name
	name := page.Find("main h1").Text()

	// Get contributors' usernames
	contributors_selection := page.Find("main .contributors a")
	contributors := make([]string, contributors_selection.Length())

	contributors_selection.Each(func(i int, link *goquery.Selection) {
		username, _ := link.Attr("data-name")
		contributors[i] = username
	})

	// Get icon svg element
	icon_target := page.Find("main svg.preview-icon").First()

	// Remove all data-* attributes from the SVG tag only
	for _, attr := range icon_target.Nodes[0].Attr {
		if strings.HasPrefix(attr.Key, "data-") {
			icon_target.RemoveAttr(attr.Key)
		}
	}

	// Get svg code
	svg, reading_err := goquery.OuterHtml(icon_target)
	if reading_err != nil {
		return icon_details{}, reading_err
	}

	return icon_details{
		Name:         name,
		Contributors: contributors,
		SVG:          svg,
	}, nil
}

func save_icon(details icon_details) error {
	// Create a new file named "myfile.txt" in the same directory
	file, err := os.Create(fmt.Sprintf("%s.svg", details.Name))
	if err != nil { // Check for an error during file creation
		return err
	}
	defer file.Close() // Ensure the file is closed when the function exits

	_, write_err := file.WriteString(details.to_string())
	if write_err != nil {
		return err
	}

	return nil
}

func main() {

	flag.Parse()

	icon_names := flag.Args()
	total_icons := len(icon_names)

	fmt.Printf("\nProcessing %d icons...\n\n", total_icons)

	state_output := func(err error) {
		if err != nil {
			fmt.Print("Failed")
			fmt.Printf("\nError: %s \n'n", err)
		} else {
			fmt.Print("Done\n")
		}
	}

	for i, icon_name := range icon_names {
		// Getting html document
		fmt.Printf("[%d/%d] Fetching Webpage...    ", i+1, total_icons)
		doc, fetching_err := fetch_page(icon_name)

		state_output(fetching_err)
		if fetching_err != nil {
			continue
		}

		// Extracting details from html document
		fmt.Printf("[%d/%d] Extracting details...  ", i+1, total_icons)
		details, extreacting_err := extract_details(doc)

		state_output(extreacting_err)
		if extreacting_err != nil {
			continue
		}

		// Save icon to file
		fmt.Printf("[%d/%d] Extracting details...  ", i+1, total_icons)
		saving_err := save_icon(details)

		state_output(saving_err)
		if saving_err != nil {
			continue
		}

		fmt.Println("")
	}
}
