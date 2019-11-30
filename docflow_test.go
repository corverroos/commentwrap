package docflow

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

)

func TestTest1(t *testing.T) {
	*limit = 40
	analyzer := Analyzer
	analyzer.RunDespiteErrors = true // importing non-standard-lib stuff errors.
	analysistest.Run(t, analysistest.TestData(), analyzer, "./...")
}

func TestFlowGroup(t *testing.T) {
	tests := []struct {
		Name string
		Input string
		Limit int
		Output string
		Flowed bool
	}{
		{
			Name: "single long line",
			Input: "This is a single long line",
			Limit: 10,
			Output:"This is a\nsingle\nlong line",
			Flowed: true,
		},{
			Name: "multi short line",
			Input: "This \nis a\n multi \nlong line",
			Limit: 10,
			Output:"This \nis a\n multi \nlong line",
			Flowed: false,
		},{
			Name: "single long sentence with short sentence",
			Input: "This is a single long line\n\nThis is ok",
			Limit: 10,
			Output:"This is a\nsingle\nlong line\n\nThis is ok",
			Flowed: true,
		},{
			Name: "single long sentence with short sentence",
			Input: "This is a single long line\n Code block ignore",
			Limit: 10,
			Output:"This is a\nsingle\nlong line\n Code block ignore",
			Flowed: true,
		},{
			Name: "sandwich long, code, code, long",
			Input: "This is a single long line\n Code\n block\n ignore\n\nThis is a single long line",
			Limit: 10,
			Output:"This is a\nsingle\nlong line\n Code\n block\n ignore\n\nThis is a\nsingle\nlong line",
			Flowed: true,
		},{
			Name: "short, long code, short",
			Input: "This is\nshort\n Long code block ignore\n\nThis is\nshort\nagain",
			Limit: 10,
			Output:"This is\nshort\n Long code block ignore\n\nThis is\nshort\nagain",
			Flowed: false,
		},{
			Name: "short, directive, short",
			Input: "This is\ngo:generate all the good stuff\n Long code block ignore\n\n",
			Limit: 10,
			Output:"This is\ngo:generate all the good stuff\n Long code block ignore\n\n",
			Flowed: false,
		},{
			Name: "short, directive, short",
			Input: "This is a long line\ngo:generate all the good stuff",
			Limit: 10,
			Output:"This is a\nlong line\ngo:generate all the good stuff",
			Flowed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res, ok := flowGroup(sp(test.Input), test.Limit)
			require.EqualValues(t, sp(test.Output), res)
			require.Equal(t, test.Flowed, ok)
		})
	}
}

func TestFlowBlock(t *testing.T) {
	tests := []struct {
		Name string
		Input string
		Limit int
		Output string
		Flowed bool
	}{
		{
			Name: "single long line",
			Input: "This is a single long line",
			Limit: 10,
			Output:"This is a\nsingle\nlong line",
			Flowed: true,
		},{
			Name: "multi short line",
			Input: "This \nis a\n multi \nlong line",
			Limit: 10,
			Output:"This \nis a\n multi \nlong line",
			Flowed: false,
		},{
			Name: "multi long line",
			Input: "This is a multi \n\nlong\nline",
			Limit: 10,
			Output:"This is a\nmulti\nlong line",
			Flowed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res, ok:= flowBlock(sp(test.Input), test.Limit)
			require.EqualValues(t, sp(test.Output), res)
			require.Equal(t, test.Flowed, ok)
		})
	}
}

func sp(s string) []string {
	return strings.Split(s, "\n")
}

