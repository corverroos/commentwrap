# docflow

docflow is a golang [analysis](https://godoc.org/golang.org/x/tools/go/analysis) providing a simple check for comments exceeding the max line limit (default 80 chars). 
It power lies in the fact that it can automatically fix it using the library [reflow](https://github.com/muesli/reflow).

## TL;DR

Install the `docflow` standalone tool.
```
go get github.com/corverroos/docflow/...
# This installs the docflow binary to $GOPATH/bin/docflow
```


Given a comment in a file `doc.go`
```
// Package foo is a great package but writing comments in golang can be frustrating
// since editing long comments
// requires manual wrapping and alignement which is tedious and wastes time that could be spent elsewhere.
package foo
```

Run the analysis with `-fix` flag.
```
docflow -fix path/to/doc.go
```

Results in:




