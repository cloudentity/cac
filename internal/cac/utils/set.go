package utils

type void struct{}

var member void

type Set[T comparable] map[T]void

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{}
}

func (s *Set[T]) Add(item T) {
	ss := *s
	ss[item] = member
}

func (s *Set[T]) Remove(item T) {
	ss := *s
	delete(ss, item)
}

func (s *Set[T]) Has(item T) bool {
	_, ok := (*s)[item]
	return ok
}
