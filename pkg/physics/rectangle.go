package physics

import r "github.com/lachee/raylib-goplus/raylib"

var _ Shape = &Rectangle{}

type Rectangle struct {
	tags []string
	Data interface{}

	r.Rectangle
}

func NewRectangle(x, y, w, h float32) *Rectangle {
	return &Rectangle{
		Rectangle: r.NewRectangle(x, y, w, h),
	}
}

func (r *Rectangle) SetData(dat interface{}) {
	r.Data = dat
}

func (r *Rectangle) GetData() interface{} {
	return r.Data
}

func (r *Rectangle) GetTags() []string {
	return r.tags
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

func (r *Rectangle) RayRec() r.Rectangle {
	return r.Rectangle
}

func (r *Rectangle) AddTags(tags ...string) {
	r.tags = append(r.tags, tags...)
}

func (r *Rectangle) RemoveTags(tags ...string) {
	for _, t := range tags {
		for i := len(r.tags) - 1; i >= 0; i-- {
			if t == r.tags[i] {
				r.tags = append(r.tags[:i], r.tags[i+1:]...)
			}
		}
	}
}

func (r *Rectangle) ClearTags() {
	r.tags = []string{}
}
