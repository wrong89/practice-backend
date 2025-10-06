package ilist

import (
	"errors"
	"slices"
)

var (
	ErrDataAlreadyExist = errors.New("data already exist")
	ErrDataNotFound     = errors.New("data not found")
	ErrInvalidID        = errors.New("invalid id")
)

type List[V any] struct {
	list []V
}

func NewList[V any]() List[V] {
	return List[V]{
		list: make([]V, 0),
	}
}

func (l *List[V]) GetLen() int {
	return len(l.list)
}

func (l *List[V]) GetData() []V {
	return l.list
}

func (l *List[V]) GetDataByID(id int) (*V, error) {
	if id < 0 {
		return nil, ErrInvalidID
	}
	if id >= len(l.list) {
		return nil, ErrDataNotFound
	}
	return &l.list[id], nil
}

func (l *List[V]) AddData(data V) (V, error) {
	l.list = append(l.list, data)
	return data, nil
}

func (l *List[V]) UpdateData(id int, updatedData V) (V, error) {
	if _, err := l.GetDataByID(id); err != nil {
		return *new(V), err
	}

	l.list[id] = updatedData
	return updatedData, nil
}

func (l *List[V]) DeleteData(id int) error {
	if _, err := l.GetDataByID(id); err != nil {
		return err
	}

	l.list = slices.Delete(l.list, id, id+1)

	return nil
}
