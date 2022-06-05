package html_parser

import "testing"

func TestHTMLParse(t *testing.T) {
	testcases := []struct {
		name   string
		html   string
		expect string
	}{
		{
			name:   "test simple <p> tag",
			html:   `<p>The way you can go</p>`,
			expect: "The way you can go\n",
		},
		{
			name:   "test multiple <p> tags",
			html:   `<p>The way you can go</p><p>The way you can go</p>`,
			expect: "The way you can go\nThe way you can go\n",
		},
		{
			name: "test multiple lines of <p> tags",
			html: `<p>The way you can go</p>
			<p>The way you can go</p>`,
			expect: "The way you can go\nThe way you can go\n",
		},
		{
			name:   "test <i> tag",
			html:   `<i>abc</i>`,
			expect: "\x1b[3ma\x1b[0m\x1b[3mb\x1b[0m\x1b[3mc\x1b[0m",
		},
		{
			name:   "test multiple <i> tags",
			html:   `<i>abc</i><i>def</i>`,
			expect: "\x1b[3ma\x1b[0m\x1b[3mb\x1b[0m\x1b[3mc\x1b[0m\x1b[3md\x1b[0m\x1b[3me\x1b[0m\x1b[3mf\x1b[0m",
		},
		{
			name:   "test <p> tag with <i> tag",
			html:   `<p>The way <i>you</i> can go</p>`,
			expect: "The way \x1b[3my\x1b[0m\x1b[3mo\x1b[0m\x1b[3mu\x1b[0m can go\n",
		},
		{
			name:   "test empty <p> tag",
			html:   `<p/>`,
			expect: "",
		},
		// {
		// 	name:   "test nested <p> tags with <i> tag(p nested is invalid)",
		// 	html:   `<p>The way <p><i>you</i></p> can go</p>`,
		// 	expect: "The way\n\x1b[3my\x1b[0m\x1b[3mo\x1b[0m\x1b[3mu\x1b[0m\ncan go\n",
		// },
	}
	for _, tc := range testcases {
		if str, err := Parse(tc.html, formater{}); err != nil {
			t.Error(err)
		} else if str.Render() != tc.expect {
			rs := str.Render()
			t.Errorf("got: %v, expect: %v", rs, tc.expect)
		}
	}

}
