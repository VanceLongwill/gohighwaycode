package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

const (
	baseURL = "https://www.gov.uk"
)

// HighwayCode represents the root node
type HighwayCode struct {
	Chapters []Chapter `xml:Chapters`
}

// Chapter represents a chapter of the highway code
type Chapter struct {
	Title    string    `xml:"Title"`
	Summary  string    `xml:"Summary"`
	URL      string    `xml:"Url"`
	Sections []Section `xml:"Sections"`
}

// Section represents a section of a chapter of the highway code
type Section struct {
	Title   string `xml:"SectionTitle"`
	Content struct {
		Inner string `xml:",cdata"`
	} `xml:"Content"`
}

func getChapters() []Chapter {
	doc, err := goquery.NewDocument(baseURL + "/guidance/the-highway-code")
	if err != nil {
		log.Fatal(err)
	}

	foundChapters := doc.Find("article#content ol.section-list li a")
	chapters := make([]Chapter, foundChapters.Length())

	foundChapters.Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		chapters[index].URL = baseURL + link

		spans := item.Find("span")
		chapters[index].Title = spans.Eq(0).Text()
		chapters[index].Summary = spans.Eq(1).Text()
	})

	return chapters
}

func getSections(link string) []Section {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		log.Fatal(err)
	}

	foundSections := doc.Find("div.gem-c-govspeak.govuk-govspeak h2")
	sections := make([]Section, foundSections.Length())

	foundSections.Each(func(index int, item *goquery.Selection) {
		sections[index].Title = item.Text()
		var untilNext *goquery.Selection
		if index == foundSections.Length()-1 {
			untilNext = item.NextAll()
		} else {
			untilNext = item.NextUntilSelection(foundSections.Eq(index + 1))
		}

		var htmlContent string
		untilNext.Each(func(ind int, selectedItem *goquery.Selection) {
			htmlSelected, _ := selectedItem.Html()
			htmlContent += htmlSelected
		})

		sections[index].Content.Inner = htmlContent
	})

	return sections
}

// Scrape gets the full highway code from the gov uk website
func Scrape() HighwayCode {
	chapters := getChapters()
	pause := time.Duration(100) * time.Millisecond

	for i := range chapters {
		chapters[i].Sections = getSections(chapters[i].URL)
		time.Sleep(pause)
	}

	return HighwayCode{
		Chapters: chapters,
	}
}
