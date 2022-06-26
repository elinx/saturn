package saturn

import (
	"reflect"
	"testing"

	"github.com/elinx/saturn/pkg/epub"
	_ "github.com/elinx/saturn/pkg/logconfig"
)

func TestParse(t *testing.T) {
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
				BlockPos: map[epub.ManifestId]int{},
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
				BlockPos: map[epub.ManifestId]int{},
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
				BlockPos: map[epub.ManifestId]int{},
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
				BlockPos: map[epub.ManifestId]int{},
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
				BlockPos: map[epub.ManifestId]int{},
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
				BlockPos: map[epub.ManifestId]int{},
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
				BlockPos: map[epub.ManifestId]int{},
			},
		},
	}
	for _, tc := range testcases {
		render := New(nil)
		if err := render.parse1(tc.html); err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(tc.expect, render.buffer) {
			t.Errorf("case %s failed: got(%d lines): \n%v\n, expect(%d lines): \n%v\n",
				tc.name,
				len(render.buffer.Lines), render.buffer,
				len(tc.expect.Lines), tc.expect)
		}
	}

}

func TestRenderWrap(t *testing.T) {
	testcases := []struct {
		name   string
		line   string
		width  int
		expect string
	}{
		{
			name:   "simple",
			line:   "The way you can go",
			width:  7,
			expect: "The way\nyou can\ngo",
		},
	}
	for _, tc := range testcases {
		render := New(nil)
		wraped, linesNum := render.renderWrap(tc.line, tc.width)
		if wraped != tc.expect {
			t.Errorf("case %s failed: got: %s, expect: %s", tc.name, wraped, tc.expect)
		}
		if linesNum != 3 {
			t.Errorf("case %s failed: got: %d, expect: %d", tc.name, linesNum, 1)
		}
	}

}
