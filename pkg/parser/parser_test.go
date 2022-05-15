package parser

import (
	"testing"

	_ "github.com/elinx/saturn/pkg/logconfig"
)

func TestP(t *testing.T) {
	testcases := []struct {
		html   string
		expect string
	}{
		{
			html:   `<p>The way you can go</p>`,
			expect: "The way you can go",
		},
		{
			html:   `<p>The way you can go</p><p>The way you can go</p>`,
			expect: "The way you can goThe way you can go",
		},
		{
			html: `<p>The way you can go</p>
<p>The way you can go</p>`,
			expect: "The way you can go\nThe way you can go",
		},
		{
			html:   `<i>The way you can go</i>`,
			expect: "\x1b[3mThe way you can go\x1b[0m",
		},
		{
			html:   `<p>The way <i>you</i> can go</p>`,
			expect: "The way \x1b[3myou\x1b[0m can go",
		},
	}
	for _, tc := range testcases {
		if str, err := Parse(tc.html, htmlFormater{}); err != nil {
			t.Error(err)
		} else if str != tc.expect {
			t.Errorf("got: %s, expect: %s", str, tc.expect)
		}
	}
}
