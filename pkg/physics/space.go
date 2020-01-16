package physics

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

type Space []Shape

// type Space []*Rectangle

func NewSpace() *Space {
	return &Space{}
}

func (s *Space) Clear() {
	*s = []Shape{}
	// *s = []*Rectangle{}
}

func (s *Space) Add(rec ...Shape) {
	*s = append(*s, rec...)
}

func (s *Space) SetData(dat interface{}) {
	for i := range *s {
		(*s)[i].SetData(dat)
	}
}

func (s *Space) RayRec() r.Rectangle {
	rec := r.Rectangle{}
	tmp := r.Rectangle{}
	set := r.Rectangle{}

	for i := range *s {
		if ((*s)[i].RayRec().X < rec.X) || (set.X == 0) {
			set.X = 1
			rec.X = (*s)[i].RayRec().X
		}
		if ((*s)[i].RayRec().Y < rec.Y) || (set.Y == 0) {
			set.Y = 1
			rec.Y = (*s)[i].RayRec().Y
		}
		if ((*s)[i].RayRec().MaxPosition().X > tmp.X) || (set.Width == 0) {
			set.Width = 1
			tmp.X = (*s)[i].RayRec().MaxPosition().X
		}
		if ((*s)[i].RayRec().MaxPosition().Y > tmp.Y) || (set.Height == 0) {
			set.Height = 1
			tmp.Y = (*s)[i].RayRec().MaxPosition().Y
		}
	}

	rec.Width = tmp.X - rec.X
	rec.Height = tmp.Y - rec.Y

	return rec
}

func (s *Space) Overlaps(rec r.Rectangle) bool {
	for i := range *s {
		if (*s)[i].Overlaps(rec) {
			return true
		}
	}

	return false
}

func (s *Space) IsColliding(s2 *Space) bool {
	for i := range *s {
		for j := range *s2 {
			if (*s)[i].Overlaps((*s2)[j].RayRec()) {
				return true
			}
		}
	}

	return false
}

func (s *Space) Remove(rec ...Shape) {
	for i := range rec {
		for j := len(*s) - 1; i >= 0; i-- {
			if rec[i] == (*s)[j] {
				*s = append((*s)[:j], (*s)[j+1:]...)
			}
		}
	}
}

func (s *Space) AddTags(tags ...string) {
	for i := range *s {
		for _, t := range tags {
			if !(*s)[i].HasTags(t) {
				(*s)[i].AddTags(t)
			}
		}
	}
}

func (s *Space) ClearTags() {
	for i := range *s {
		(*s)[i].ClearTags()
	}
}

func (s *Space) GetTags() []string {
	var tmp = &Rectangle{}
	for i := range *s {
		tt := (*s)[i].GetTags()
		for _, t := range tt {
			if !tmp.HasTags(t) {
				tmp.AddTags(t)
			}
		}
	}

	return tmp.GetTags()
}

func (s *Space) HasTags(tags ...string) bool {
	for i := range *s {
		for _, t := range tags {
			if (*s)[i].HasTags(t) {
				return true
			}
		}
	}
	return false
}

func (s *Space) RemoveTags(tags ...string) {
	for i := range *s {
		(*s)[i].RemoveTags(tags...)
	}
}

func (s *Space) Position() r.Vector2 {
	var pos r.Vector2
	var set r.Vector2

	for i := range *s {
		if (*s)[i].Position().X < pos.X || set.X == 0 {
			pos.X = (*s)[i].Position().X
			set.X = 1
		}

		if (*s)[i].Position().Y < pos.X || set.Y == 0 {
			pos.Y = (*s)[i].Position().Y
			set.Y = 1
		}
	}

	return pos
}

func (s *Space) GetData() interface{} {
	if len(*s) > 0 {
		return (*s)[0].GetData()
	}

	return nil
}

func (s *Space) Filter(filter func(Shape) bool) *Space {
	subSpace := &Space{}
	for i := range *s {
		if filter((*s)[i]) {
			subSpace.Add((*s)[i])
		}
	}

	return subSpace
}

func (s *Space) FilterByTags(tags ...string) *Space {
	return s.Filter(func(r Shape) bool {
		if r.HasTags(tags...) {
			return true
		}
		return false
	})
}

func (s *Space) FilterOutByTags(tags ...string) *Space {
	return s.Filter(func(r Shape) bool {
		if r.HasTags(tags...) {
			return false
		}
		return true
	})
}
