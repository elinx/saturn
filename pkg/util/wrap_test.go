package util

import (
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	tt := []struct {
		Input         string
		Expected      string
		Limit         int
		KeepNewlines  bool
		PreserveSpace bool
		TabWidth      int
	}{
		// No-op, should pass through, including trailing whitespace:
		{
			Input:         "foobar\n ",
			Expected:      "foobar\n ",
			Limit:         0,
			KeepNewlines:  true,
			PreserveSpace: false,
			TabWidth:      0,
		},
		// Nothing to wrap here, should pass through:
		{
			Input:         "foo",
			Expected:      "foo",
			Limit:         4,
			KeepNewlines:  true,
			PreserveSpace: false,
			TabWidth:      0,
		},
		// In contrast to wordwrap we break a long word to obey the given limit
		{
			Input:         "foobarfoo",
			Expected:      "foob\narfo\no",
			Limit:         4,
			KeepNewlines:  true,
			PreserveSpace: false,
			TabWidth:      0,
		},
		// Newlines in the input are respected if desired
		{
			Input:         "f\no\nobar",
			Expected:      "f\no\noba\nr",
			Limit:         3,
			KeepNewlines:  true,
			PreserveSpace: false,
			TabWidth:      0,
		},
		// Leading whitespaces after forceful line break can be preserved if desired
		{
			Input:         "foo bar\n  baz",
			Expected:      "foo\n ba\nr\n  b\naz",
			Limit:         3,
			KeepNewlines:  true,
			PreserveSpace: true,
			TabWidth:      0,
		},
		// Tabs are broken up according to the configured TabWidth
		// {
		// 	Input:         "foo\tbar",
		// 	Expected:      "foo\t\nbar",
		// 	Limit:         4,
		// 	KeepNewlines:  true,
		// 	PreserveSpace: true,
		// 	TabWidth:      3,
		// },
		{
			Input:    "道可道非常道",
			Expected: "道\n可\n道\n非\n常\n道",
			Limit:    2,
		},
		{
			Input:    "道可道非常道",
			Expected: "道\n可\n道\n非\n常\n道",
			Limit:    3,
		},
		{
			Input:    "道可道非常道",
			Expected: "道可\n道非\n常道",
			Limit:    4,
		},
		{
			Input:    "道可道非常道",
			Expected: "道可\n道非\n常道",
			Limit:    5,
		},
		{
			Input:    "道可道非常道",
			Expected: "道可道\n非常道",
			Limit:    6,
		},
		{
			Input:    "道可道非常道",
			Expected: "道可道非\n常道",
			Limit:    9,
		},
		{
			Input:    "道可道非常道",
			Expected: "道可道非常道",
			Limit:    13,
		},
		{
			Input:    "道tao可道非常道",
			Expected: "道\nta\no\n可\n道\n非\n常\n道",
			Limit:    2,
		},
		{
			Input:    "道tao可道非常道",
			Expected: "道t\nao\n可\n道\n非\n常\n道",
			Limit:    3,
		},
		{
			Input:    "道\x1b[7m可道非常\x1b[0m道",
			Expected: "道\x1b[7m可道\n非常\x1b[0m道",
			Limit:    6,
		},
		{
			Input:    "道\x1b[7m可道非常\x1b[0m道",
			Expected: "道\x1b[7m可道非\n常\x1b[0m道",
			Limit:    9,
		},
		{
			Input:    "道\x1b[7m可道非常\x1b[0m道",
			Expected: "道\x1b[7m可道非常\x1b[0m道",
			Limit:    13,
		},
		// {
		// 	Input:    "\x1b[7m道\x1b[0m\x1b[7m可\x1b[0m\x1b[7m道\x1b[0m\x1b[7m非\x1b[0m\x1b[7m常\x1b[0m\x1b[7m道\x1b[0m",
		// 	Expected: "\x1b[7m道\x1b[0m\x1b[7m可\x1b[0m\x1b[7m道\x1b[0m\n\x1b[7m非\x1b[0m\x1b[7m常\x1b[0m\x1b[7m道\x1b[0m",
		// 	Limit:    6,
		// },
	}

	for i, tc := range tt {
		actual := Wrap(tc.Input, tc.Limit)
		if actual != tc.Expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.Expected, actual)
		}
	}
}

func TestWrapPos(t *testing.T) {
	tt := []struct {
		Input    string
		Limit    int
		Expected int
		VX, VY   int
	}{
		{
			Input:    "foobar",
			Limit:    4,
			Expected: 0,
			VX:       0,
			VY:       0,
		},
		{
			Input:    "foobar",
			Limit:    4,
			Expected: 5,
			VX:       1,
			VY:       1,
		},
		{
			Input:    "foo\nbar",
			Limit:    4,
			Expected: 6,
			VX:       2,
			VY:       1,
		},
		{
			Input:    "foo\n\nbar",
			Limit:    4,
			Expected: 7,
			VX:       2,
			VY:       2,
		},
		{
			Input:    "道可道非常道",
			Limit:    6,
			Expected: 0,
			VX:       0,
			VY:       0,
		},
		{
			Input:    "道可道非常道",
			Limit:    6,
			Expected: 0,
			VX:       1,
			VY:       0,
		},
		{
			Input:    "道可道非常道",
			Limit:    12,
			Expected: 1,
			VX:       2,
			VY:       0,
		},
		{
			Input:    "道可道非常道",
			Limit:    12,
			Expected: 1,
			VX:       3,
			VY:       0,
		},
		{
			Input:    "道可道非常道",
			Limit:    12,
			Expected: 2,
			VX:       4,
			VY:       0,
		},
		{
			Input:    "道可道非常道",
			Limit:    12,
			Expected: 4,
			VX:       9,
			VY:       0,
		},
		{
			Input:    "道可道非常道",
			Limit:    6,
			Expected: 3,
			VX:       1,
			VY:       1,
		},
		{
			Input:    "道\x1b[7m可道非常\x1b[0m道",
			Limit:    6,
			Expected: 3,
			VX:       1,
			VY:       1,
		},
	}

	for i, tc := range tt {
		actual := LocBeforeWraped(tc.Input, tc.Limit, tc.VX, tc.VY)
		wraps := strings.Split(Wrap(tc.Input, tc.Limit), "\n")
		if actual != tc.Expected {
			t.Errorf("Test %d, expected:\n\n`%v:%c`\n\nActual Output:\n\n`%v:%c`",
				i, tc.Expected, wraps[tc.VY][tc.VX], actual, tc.Input[actual])
		}
	}
}
