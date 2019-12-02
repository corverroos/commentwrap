// Command commentwrap runs the comment wrap analyzer.
package main

import (
	"github.com/corverroos/commentwrap"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(commentwrap.Analyzer)
}
