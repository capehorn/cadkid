package lang

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(data T) {
	s.items = append(s.items, data)
}

func (s *Stack[T]) Pop() {
	if s.IsEmpty() {
		return
	}
	s.items = s.items[:len(s.items)-1]
}

func (s *Stack[T]) Top() (any, error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}
