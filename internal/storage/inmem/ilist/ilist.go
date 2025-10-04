package ilist

import (
	"errors"
)

var (
	ErrDataAlreadyExist = errors.New("data already exist")
	ErrDataNotFound     = errors.New("data not found")
)

type List[V any] struct {
	list map[int]V
	// in use last id
	lastID int
}

func NewList[V any]() List[V] {
	return List[V]{
		list:   make(map[int]V),
		lastID: 0,
	}
}

func (l *List[V]) GetLastID() int {
	return l.lastID
}

func (l *List[V]) GetData(id int) (V, error) {
	result, ok := l.list[id]
	if !ok {
		return *new(V), ErrDataNotFound
	}

	return result, nil
}

func (l *List[V]) AddData(data V) (V, error) {
	if _, err := l.GetData(l.lastID + 1); err == nil {
		return *new(V), ErrDataAlreadyExist
	}
	l.lastID++
	l.list[l.lastID] = data

	return data, nil
}

func (l *List[V]) UpdateData(id int, updatedData V) (V, error) {
	if _, ok := l.list[id]; !ok {
		return *new(V), ErrDataNotFound
	}

	l.list[id] = updatedData
	return updatedData, nil
}

func (l *List[V]) DeleteData(id int) error {
	if _, err := l.GetData(id); err != nil {
		return err
	}

	delete(l.list, id)
	return nil
}
