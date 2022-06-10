package parser

import (
	"reflect"
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
		{
			html:   `<p/>`,
			expect: "",
		},
	}
	for _, tc := range testcases {
		if str, err := Parse(tc.html, HtmlFormater{}); err != nil {
			t.Error(err)
		} else if str != tc.expect {
			t.Errorf("got: %s, expect: %s", str, tc.expect)
		}
	}
}

func TestRenderParse(t *testing.T) {
	testcases := []struct {
		name   string
		html   string
		expect *Buffer
	}{
		{
			name: "simple",
			html: `<p>The way you can go</p>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way you can go", Style: "", Pos: 0},
						},
						Style: "p",
					},
				},
			},
		},
		{
			name: "empty p",
			html: `<p/>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content:  "",
						Segments: nil,
						Style:    "p",
					},
				},
			},
		},
		{
			name: "two p",
			html: `<p>The way you can go</p><p>The way you can go</p>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way you can go", Style: "", Pos: 0},
						},
						Style: "p",
					},
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way you can go", Style: "", Pos: 0},
						},
						Style: "p",
					},
				},
			},
		},
		{
			name: "two p two line with white spaces",
			html: `<p>The way you can go</p>

			<p>The way you can go</p>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way you can go", Style: "", Pos: 0},
						},
						Style: "p",
					},
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way you can go", Style: "", Pos: 0},
						},
						Style: "p",
					},
				},
			},
		},
		{
			name: "simple italic",
			html: `<i>The way you can go</i>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way you can go", Style: "", Pos: 0},
						},
						Style: "i",
					},
				},
			},
		},
		{
			name: "p with nested italic",
			html: `<p>The way <i>you</i> can go</p>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way you can go",
						Segments: []Segment{
							{Content: "The way ", Style: "", Pos: 0},
							{Content: "you", Style: "i", Pos: 8},
							{Content: " can go", Style: "", Pos: 11},
						},
						Style: "p",
					},
				},
			},
		},
		{
			name: "2p with italic in between",
			html: `<p>The way</p><i>you</i><p>can go</p>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way",
						Segments: []Segment{
							{Content: "The way", Style: "", Pos: 0},
						},
						Style: "p",
					},
					{
						Content: "you",
						Segments: []Segment{
							{Content: "you", Style: "", Pos: 0},
						},
						Style: "i",
					},
					{
						Content: "can go",
						Segments: []Segment{
							{Content: "can go", Style: "", Pos: 0},
						},
						Style: "p",
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		render := New(nil)
		if err := render.Parse(tc.html); err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(tc.expect, render.buffer) {
			t.Errorf("case %s failed: got(%d lines): \n%v\n, expect(%d lines): \n%v\n",
				tc.name,
				len(render.buffer.Lines), render.buffer,
				len(tc.expect.Lines), tc.expect)
		}
	}

}
