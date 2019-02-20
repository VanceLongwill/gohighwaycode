package main

import (
	"bytes"
	"fmt"
	"github.com/bmaupin/go-epub"
	"html/template"
	"io"
	"log"
	"strings"
)

type templateHelper struct {
	chapterTmpl  *template.Template
	contentsTmpl *template.Template
}

func newTemplateHelper() *templateHelper {
	chapterTmplFile := "assets/templates/chapter.html"
	contentsTmplFile := "assets/templates/contents.html"
	chapterTmpl := template.Must(template.ParseFiles(chapterTmplFile))
	contentsTmpl := template.Must(template.ParseFiles(contentsTmplFile))
	return &templateHelper{
		chapterTmpl:  chapterTmpl,
		contentsTmpl: contentsTmpl,
	}
}

type chapterCover struct {
	Title   string
	Summary string
}
type contentsInfo struct {
	Chapters []Chapter
}

func (t *templateHelper) chapter(title, summary string) (string, error) {
	var buf bytes.Buffer
	wr := io.Writer(&buf)
	data := chapterCover{
		Title:   title,
		Summary: summary,
	}
	if err := t.chapterTmpl.Execute(wr, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t *templateHelper) contents(hwCode HighwayCode) (string, error) {
	var buf bytes.Buffer
	wr := io.Writer(&buf)
	if err := t.contentsTmpl.Execute(wr, hwCode); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Generate creates an epub
func Generate(hwCode HighwayCode) *epub.Epub {
	tmplHelper := newTemplateHelper()
	chapters := hwCode.Chapters

	e := epub.NewEpub("The Highway Code")
	e.SetAuthor("GOV.UK")
	e.SetLang("English")

	css, cssErr := e.AddCSS("assets/css/style.css", "style.css")
	if cssErr != nil {
		log.Fatal(cssErr)
	}

	contentsPage, contentsErr := tmplHelper.contents(hwCode)
	if contentsErr != nil {
		log.Fatal(contentsErr)
	}
	_, err := e.AddSection(contentsPage, "Contents", "", css)
	if err != nil {
		log.Fatal(err)
	}

	for index, chapter := range chapters {
		chapterTitle := fmt.Sprintf("CHAPTER %d: %s", index+1, strings.ToUpper(chapter.Title))
		chapterPage, err := tmplHelper.chapter(chapterTitle, chapter.Summary)
		if err != nil {
			log.Fatal(err)
		}
		_, chapterErr := e.AddSection(chapterPage, chapterTitle, "", css)
		if chapterErr != nil {
			log.Fatal(chapterErr)
		}
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
