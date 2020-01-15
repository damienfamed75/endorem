package physics

import r "github.com/lachee/raylib-goplus/raylib"

var _ Shape = &Rectangle{}

type Rectangle struct {
	tags []string

	r.Rectangle
}

func NewRectangle(x, y, w, h float32) *Rectangle {
	return &Rectangle{
		Rectangle: r.NewRectangle(x, y, w, h),
	}
}

func (r *Rectangle) HasTags(tags ...string) bool {
	hasTags := true

	for _, wanted := range tags {
		found := false
		for _, shapeTag := range r.tags {
			if wanted == shapeTag {
				found = true
				continue
			}
		}
		if !found {
			hasTags = false
			break
		}
	}

	return hasTags
}

func (r *Rectangle) AddTags(tags ...string) {
	r.tags = append(r.tags, tags...)
}

func (r *Rectangle) ClearTags() {
	r.tags = []string{}
}
