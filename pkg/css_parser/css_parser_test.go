package cssparser

import (
	"reflect"
	"testing"
)

func TestCssTokens(t *testing.T) {
	testcases := []struct {
		css    string
		expect []tokenInfo
	}{
		{
			css: `p { color: red; }`,
			expect: []tokenInfo{
				{IdentifierToken, "p", 0, 1},
				{OpenCurlyBraceToken, "{", 2, 1},
				{IdentifierToken, "color", 4, 5},
				{ColonToken, ":", 9, 1},
				{IdentifierToken, "red", 11, 3},
				{SemicolonToken, ";", 14, 1},
				{CloseCurlyBraceToken, "}", 16, 1},
			},
		},
		{
			css: `p { border-top: 1px solid black; }`,
			expect: []tokenInfo{
				{IdentifierToken, "p", 0, 1},
				{OpenCurlyBraceToken, "{", 2, 1},
				{IdentifierToken, "border-top", 4, 10},
				{ColonToken, ":", 14, 1},
				{IdentifierToken, "1px", 16, 3},
				{IdentifierToken, "solid", 20, 5},
				{IdentifierToken, "black", 26, 5},
				{SemicolonToken, ";", 31, 1},
				{CloseCurlyBraceToken, "}", 33, 1},
			},
		},
		{
			css: `p { color: red; }
				.p { color: red; }`,
			expect: []tokenInfo{
				{IdentifierToken, "p", 0, 1},
				{OpenCurlyBraceToken, "{", 2, 1},
				{IdentifierToken, "color", 4, 5},
				{ColonToken, ":", 9, 1},
				{IdentifierToken, "red", 11, 3},
				{SemicolonToken, ";", 14, 1},
				{CloseCurlyBraceToken, "}", 16, 1},
				{IdentifierToken, ".p", 22, 2},
				{OpenCurlyBraceToken, "{", 25, 1},
				{IdentifierToken, "color", 27, 5},
				{ColonToken, ":", 32, 1},
				{IdentifierToken, "red", 34, 3},
				{SemicolonToken, ";", 37, 1},
				{CloseCurlyBraceToken, "}", 39, 1},
			},
		},
	}
	for _, tc := range testcases {
		tokenizer := NewTokenizer(tc.css)
		tokenizer.lex()
		actual := tokenizer.GetTokens()
		if !reflect.DeepEqual(actual, tc.expect) {
			t.Errorf("got: %v, expect: %v", actual, tc.expect)
		}
	}
}

func TestRules(t *testing.T) {
	testcases := []struct {
		css    string
		expect []*rule
	}{
		{
			css: `p { color: red; }`,
			expect: []*rule{
				{
					selector: "p",
					declarations: []declaration{
						{
							property: "color",
							value:    "red",
						},
					},
				},
			},
		},
		{
			css: `p { border-top: 1px solid black; }`,
			expect: []*rule{
				{
					selector: "p",
					declarations: []declaration{
						{
							property: "border-top",
							value:    "1px solid black",
						},
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		parser := NewParser()
		rules, err := parser.Parse(tc.css)
		if err != nil {
			t.Errorf("failed to parse css: %v", err)
		}
		if !reflect.DeepEqual(rules, tc.expect) {
			t.Errorf("got: %v, expect: %v", rules, tc.expect)
		}
	}
}
