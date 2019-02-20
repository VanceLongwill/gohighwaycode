package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	outputFmt := flag.String("format", "", "generates a highway code book in specified the output format")
	shouldUpdate := flag.Bool("update", false, "update fetches the latest highway code content")
	xmlFn := flag.String("source", "highwayCode.xml", "specify a filename for the xml data source")
	flag.Parse()

	var shouldGenerate bool
	switch *outputFmt {
	case "epub":
		shouldGenerate = true
	case "":
		shouldGenerate = false
	default:
		fmt.Println("Invalid output format")
		shouldGenerate = false
	}

	if !shouldGenerate && !*shouldUpdate {
		flag.Usage()
		return
	}

	var hwCode HighwayCode
	if *shouldUpdate {
		fmt.Print("Fetching highway code data from the gov uk site...")
		hwCode = Scrape()
		fmt.Println("Done")

		highwayCodeXML, xmlErr := xml.MarshalIndent(hwCode, "", "  ")
		if xmlErr != nil {
			log.Fatal(xmlErr)
		}

		fmt.Printf("Saving highway code data to: %s ...", *xmlFn)
		if err := ioutil.WriteFile(*xmlFn, highwayCodeXML, 0644); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Done")
		}
	} else {
		fmt.Printf("Fetching highway code data from file: %s ...", *xmlFn)
		rawXML, readErr := ioutil.ReadFile(*xmlFn)
		if readErr != nil {
			fmt.Println("Please run with -update flag to fetch the highway code data")
			log.Fatal(readErr)
		} else if len(rawXML) < 10 {
			fmt.Println("Please run with -update flag to fetch the highway code data")
		} else if err := xml.Unmarshal(rawXML, &hwCode); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Done")
		}
	}

	if shouldGenerate {
		fmt.Printf("Generating the highway code book in %s format ...", *outputFmt)
		epub := Generate(hwCode)
		epub.Write("highwaycode.epub")
	}

	fmt.Println("Done")
}
