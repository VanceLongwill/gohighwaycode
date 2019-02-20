# UK highway code EPUB generator
> source: https://www.gov.uk/guidance/the-highway-code

CLI tool for generating an EPUB book of the highway code, lifted verbatum from the gov uk website.

## Downloads

[Click here](downloads/highwaycode.epub) to download the latest highway code EPUB

## Installation

Either:

- Download and unzip the executable [here](/downloads/gohighwaycode.zip), or
- Build & install from source by cloning this repo and running `go install`

## Usage

```
Usage of gohighwaycode:
  -format string
    	generates a highway code book in specified the output format
  -source string
    	specify a filename for the xml data source (default "highwayCode.xml")
  -update
    	update fetches the latest highway code content
```

#### Output formats
- epub (currently the only format available)

> NB: to build the ebook for the first time you must use the -update option

Running the program with the 
```
gohighwaycode -format epub
```
will output a `highwaycode.epub` file in the current directory

## Credits 


- Cross compatible EPUB css by *Matt Harrison*
    + source: https://github.com/mattharrison/epub-css-starter-kit
