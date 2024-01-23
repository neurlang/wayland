package main

import (
	"flag"
)

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "i", "", "input xml file name")

	flag.Parse()

	xmlProcess(inputFile, inputFile+".go")
}
