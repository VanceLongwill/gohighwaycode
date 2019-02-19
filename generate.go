package main

import (
	"fmt"
	"github.com/bmaupin/go-epub"
	"log"
)

// Generate creates an epub
func Generate(hwCode HighwayCode) *epub.Epub {
	chapters := hwCode.Chapters
	e := epub.NewEpub("The Highway Code")
	e.SetAuthor("GOV.UK")
	e.SetLang("English")

	css, cssErr := e.AddCSS("assets/css/style.css", "style.css")
	if cssErr != nil {
		log.Fatal(cssErr)
	}

	for index, chapter := range chapters {
		for subIndex, section := range chapter.Sections {
			title := fmt.Sprintf("%d.%d %s", index+1, subIndex+1, section.Title)
			_, err := e.AddSection(section.Content.Inner, title, "", css)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return e
}
