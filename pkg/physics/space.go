package physics

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

func (s *Space) AddTags(tags ...string) {
	for _, sh := range *s {
		sh.AddTags(tags...)
	}
}

func (s *Space) Move(x, y float32) {

}

func (s *Space) Filter(filter func(Shape) bool) *Space {
	subSpace := &Space{}
	for _, rec := range *s {
		if filter(rec) {
			subSpace.Add(rec)
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
