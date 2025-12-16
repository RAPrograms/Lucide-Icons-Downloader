package main

import (
	"errors"
	"fmt"
	"net/http"
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

func main() {
	doc, err := fetch_page("Download")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(doc.Find("main h1").Text())
}
