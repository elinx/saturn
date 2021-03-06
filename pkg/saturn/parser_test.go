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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
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
				BlockPos: map[epub.ManifestId]BufferLineIndex{},
			},
		},
		{
			name: "url embeded in p",
			html: `<p>The way you can go <a href="http://www.google.com">google</a></p>`,
			expect: &Buffer{
				Lines: []Line{
					{
						Content: "The way you can go google",
						Segments: []Segment{
							{Content: "The way you can go ", Style: "", Pos: 0},
							{Content: "google", Style: "a", Pos: 19},
						},
						Style: "p",
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		render := NewParser(nil)
		if err := render.parse1(tc.html); err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(tc.expect.Lines, render.buffer.Lines) {
			t.Errorf("case %s failed: got(%d lines): \n%v\n, expect(%d lines): \n%v\n",
				tc.name,
				len(render.buffer.Lines), render.buffer,
				len(tc.expect.Lines), tc.expect)
		}
	}

}
