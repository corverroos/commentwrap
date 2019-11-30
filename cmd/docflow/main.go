package main

import (
	"github.com/corverroos/docflow"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(docflow.Analyzer)
}
