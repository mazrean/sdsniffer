package main

import (
	"github.com/mazrean/sdsniffer"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(sdsniffer.Analyzer)
}
