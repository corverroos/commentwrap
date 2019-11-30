package main

import (
	"github.com/corverroos/docfmt"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(docflow.Analyzer)
}
