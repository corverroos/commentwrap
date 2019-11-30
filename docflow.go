package docflow

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"strings"

	"github.com/muesli/reflow"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "docflow",
	Doc: "docflow checks max line length of godoc sentence comments. " +
		"It supports fixing with -fix flag (requires subsequent go fmt in some cases).",
	Run: run,
}

var limit = flag.Int("docflow_limit", 80, "docflow line length limit")

func run(p *analysis.Pass) (interface{}, error) {
	for _, file := range p.Files {
		for _, group := range file.Comments {
			if !mustFlow(group, *limit) {
				continue
			}

			lines := strings.Split(group.Text(), "\n")
			lines = lines[:len(lines)-1]
			text, flowed := flowGroup(lines, *limit)
			if !flowed {
				continue
			}

			pos := p.Fset.Position(group.Pos())
			indent := strings.Repeat("\t", pos.Column-1)
			var buf bytes.Buffer
			for i, line := range text {
				tabs := indent
				prefix := "//"
				if i == 0 {
					tabs = ""
					prefix = "/" // Workaround for fixes on pos 0.
				}
				if !isDirective(line) {
					prefix += " "
				}
				suffix := "\n"
				if len(text) == i+1 {
					suffix = ""
				}
				buf.WriteString(tabs + prefix + line + suffix)
			}

			p.Report(analysis.Diagnostic{
				Pos:     group.Pos(),
				End:     group.End(),
				Message: fmt.Sprintf("Comment exceeds %d character limit", *limit),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: "Reflow",
					TextEdits: []analysis.TextEdit{{
						Pos:     group.Pos() + 1, // Workaround for fixes on pos 0.
						End:     group.End(),
						NewText: buf.Bytes(),
					}},
				}},
			})

		}
	}

	return nil, nil
}

// mustFlow returns true if the comment group has a line exceeding the limit.
// "/*" type comments are not supported.
func mustFlow(g *ast.CommentGroup, limit int) bool {
	var exceed bool
	// "/*" comment blocks not supported.
	for _, comment := range g.List {
		if strings.HasPrefix(comment.Text, "/*") {
			return false
		}
		if len(comment.Text) > limit {
			exceed = true
		}

	}
	return exceed
}

// flowGroup will returns true and the comment group with godoc sentences "reflowed"
// if they contain lines longer than limit.
func flowGroup(group []string, limit int) ([]string, bool) {
	const (
		normal    = 0
		indented  = 1
		empty     = 2
		directive = 3
	)
	var (
		flowed    bool
		lines     []string
		block     []string
		blockType int
	)

	flushBlock := func(block []string, typ int) {
		if typ == normal {
			// Maybe flow normal block
			var ok bool
			block, ok = flowBlock(block, limit)
			if ok {
				flowed = ok
			}
		}
		lines = append(lines, block...)
	}

	for _, line := range group {
		var typ int
		if line == "" {
			typ = empty
		} else if isDirective(line) {
			typ = directive
		} else if line[0] == ' ' || line[0] == '\t' {
			typ = indented
		} else {
			typ = normal
		}
		if typ != blockType {
			// New type, flush block
			flushBlock(block, blockType)
			block = nil
		}
		block = append(block, line)
		blockType = typ
	}

	flushBlock(block, blockType)

	return lines, flowed
}

// isDirective returns true if the unescaped comment is a go directive.
// See https://golang.org/cmd/compile/#hdr-Compiler_Directives.
func isDirective(line string) bool {
	return strings.HasPrefix(line, "go:")
}

// flowBlock will "reflow" the whole block if any line is longer than limit
// or it will return the block as is.
func flowBlock(block []string, limit int) ([]string, bool) {
	if len(block) == 0 {
		return nil, false
	}
	var overflow bool
	for _, line := range block {
		if len(line) > limit {
			overflow = true
			break
		}
	}
	if !overflow {
		return block, false
	}
	f := reflow.NewReflow(limit)
	for _, line := range block {
		f.Write([]byte(line))
		f.Write([]byte(" "))
	}
	return strings.Split(f.String(), "\n"), true
}
