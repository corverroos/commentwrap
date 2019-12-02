# commentwrap

Commentwrap is a golang [analysis](https://godoc.org/golang.org/x/tools/go/analysis) providing a check for comments exceeding the max line limit (default 80 chars). 
Its power lies in the fact that it can automatically fix it using the library [reflow](https://github.com/muesli/reflow).

## TL;DR

Install the `commentwrap` binary.
```
go get github.com/corverroos/commentwrap/...
```

Given a file `doc.go`
```
// Lorem ipsum dolor sit amet, 
// cum ea propriae lobortis reprimique, sed dolorum cotidieque ne, quo ad esse error. 
// Tempor petentium ad per, in alii detracto reprehendunt 
// mei, utamur vivendo vim ut.
//   fmt.Println("Code blocks are not wrapped")
package lipsum
```

Run the analysis with `-fix` flag.
```
$GOPATH/bin/commentwrap -fix path/to/doc.go
```

Results in:
```
// Lorem ipsum dolor sit amet, cum ea propriae lobortis reprimique, sed dolorum
// cotidieque ne, quo ad esse error. Tempor petentium ad per, in alii detracto
// reprehendunt mei, utamur vivendo vim ut.
//   fmt.Println("Code blocks are not wrapped")
package lipsum
```

# Notes

- As an analysis, its usage it the same as `go vet`. See this great [video](https://www.youtube.com/watch?v=10IMWTpCSIQ) on writing a go analysis.
- This analysis is aimed at formatting multi-line comment paragraphs.
- It is triggered by a single line in a paragraph exceeding the character limit (defaults to 80).
- Paragraphs where no lines exceed the limit are ignored.
- Comment blocks using `/* ... */` are not supported.
- It aims to support [godoc](https://blog.golang.org/godoc-documenting-go-code) and other golang comment artifacts:
  - Pre-formatted code (indented lines) are ignored.
  - [Directives](https://golang.org/cmd/compile/#hdr-Compiler_Directives) are ignored; `//go:generate`.
  - Notes are ignored; `//TODO`, `//BUG`, `//FIXME`.
- Using it as a [Goland external tool](goland.jpg) is simple.
