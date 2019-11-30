package a

import "fmt"

//go:generate ignore this very long line since it is go command

/* Star comments are also ignored since they are weird
and I'm not sure when people actually use them*/

// Exported is an exported function with a long godoc comment. // want "Comment exceeds 40 character limit"
func Exported(/* sneak comment to be ignored since it is starred */) {
	var i int

	// This is a short comment
	if i == 0 {
		// Inline very long comment that must be flowed. // want "Comment exceeds 40 character limit"
		fmt.Print("Gotcha")
	}
}
