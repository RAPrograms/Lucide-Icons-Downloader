package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type icon_details struct {
	Name         string
	Contributors []string
	SVG          string
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
	fmt.Println("Name: ", name)

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

func main() {
	doc, err := fetch_page("Download")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(extract_details(doc))
}
