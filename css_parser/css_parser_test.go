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
				{IdentifierToken, "p"},
				{OpenCurlyBraceToken, "{"},
				{IdentifierToken, "color"},
				{ColonToken, ":"},
				{IdentifierToken, "red"},
				{SemicolonToken, ";"},
				{CloseCurlyBraceToken, "}"},
			},
		},
		{
			css: `p { color: red; }
				.p { color: red; }`,
			expect: []tokenInfo{
				{IdentifierToken, "p"},
				{OpenCurlyBraceToken, "{"},
				{IdentifierToken, "color"},
				{ColonToken, ":"},
				{IdentifierToken, "red"},
				{SemicolonToken, ";"},
				{CloseCurlyBraceToken, "}"},
				{IdentifierToken, ".p"},
				{OpenCurlyBraceToken, "{"},
				{IdentifierToken, "color"},
				{ColonToken, ":"},
				{IdentifierToken, "red"},
				{SemicolonToken, ";"},
				{CloseCurlyBraceToken, "}"},
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
