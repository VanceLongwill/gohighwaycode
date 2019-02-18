package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"time"
)

const (
	baseURL = "https://www.gov.uk/guidance/the-highway-code"
)

type chapter struct {
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	URL      string    `json:"url"`
	Sections []section `json:"sections"`
}

type section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func getChapters() []chapter {
	doc, err := goquery.NewDocument(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	foundChapters := doc.Find("article#content ol.section-list li a")
	chapters := make([]chapter, foundChapters.Length())

	foundChapters.Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		chapters[index].URL = baseURL + link

		spans := item.Find("span")
		chapters[index].Title = spans.First().Text()
		chapters[index].Summary = spans.Last().Text()
	})

	return chapters
}

func getSections(link string) []section {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		log.Fatal(err)
	}

	foundSections := doc.Find("div.gem-c-govspeak.govuk-govspeak h2")
	sections := make([]section, foundSections.Length())

	foundSections.Each(func(index int, item *goquery.Selection) {
		sections[index].Title = item.Text()
		var untilNext *goquery.Selection
		if index == foundSections.Length()-1 {
			untilNext = item.NextAll()
		} else {
			untilNext = item.NextUntilSelection(foundSections.Eq(index + 1))
		}

		htmlContent, _ := untilNext.Html()
		sections[index].Content = htmlContent
	})

	return sections
}

func main() {
	fmt.Println("Starting scrape")

	chapters := getChapters()
	pause := time.Duration(100) * time.Millisecond

	for _, chapter := range chapters {
		chapter.Sections = getSections(chapter.URL)
		time.Sleep(pause)
	}

	json, jsonErr := json.Marshal(chapters)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fn := "highwayCode.json"
	if err := ioutil.WriteFile(fn, json, 0644); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Done.")
	}
}
