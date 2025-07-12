package lang

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(data T) {
	s.items = append(s.items, data)
}

func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var absent T
		return absent, fmt.Errorf("stack is empty")
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, nil
}

func (s *Stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
		var absent T
		return absent, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}

func (s *Stack[T]) Length() int {
	return len(s.items)
}
