package db

type AnnotationType string

const (
	AnnotationHighlight  AnnotationType = "Highlight"
	AnnotationUnderscore AnnotationType = "Underscore"
	AnnotationComment    AnnotationType = "Comment"
)

type AnnotationOption func(*Annotation)

type Annotation struct {
	Type    AnnotationType `json:"type"`
	Text    string         `json:"text"`
	StartX  int            `json:"startx"`
	StartY  int            `json:"starty"`
	EndX    int            `json:"endx"`
	EndY    int            `json:"endy"`
	Color   string         `json:"color"`
	Author  string         `json:"author"`
	Date    string         `json:"date"`
	Comment string         `json:"comment"`
}

func NewAnnotation(t AnnotationType, text string, opts ...AnnotationOption) Annotation {
	anno := Annotation{Type: t, Text: text}
	for _, opt := range opts {
		opt(&anno)
	}
	return anno
}

func WithComment(comment string) AnnotationOption {
	return func(a *Annotation) {
		a.Comment = comment
	}
}

func WithLocation(startx, starty, endx, endy int) AnnotationOption {
	return func(a *Annotation) {
		a.StartX = startx
		a.StartY = starty
		a.EndX = endx
		a.EndY = endy
	}
}
