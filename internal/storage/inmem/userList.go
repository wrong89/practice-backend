package inmem

import (
	"context"
	"errors"
	"practice-backend/internal/models/user"
	"practice-backend/internal/storage/inmem/ilist"
	"sync"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// Concurrent-Use
type UserList struct {
	list ilist.List[user.User]
	mtx  *sync.Mutex
}

func NewUserList() UserList {
	return UserList{
		list: ilist.NewList[user.User](),
		mtx:  new(sync.Mutex),
	}
}

func (ul *UserList) GetUserByID(ctx context.Context, id int) (user.User, error) {
	ul.mtx.Lock()
	defer ul.mtx.Unlock()

	u, err := ul.list.GetDataByID(id)
	if err != nil {
		if errors.Is(err, ilist.ErrDataNotFound) {
			return user.User{}, ErrUserNotFound
		}
		return user.User{}, err
	}

	return *u, nil
}

func (el *UserList) CreateUser(
	ctx context.Context,
	login string,
	password string,
	name string,
	surname string,
	patronymic string,
	phone string,
	email string,
) (user.User, error) {
	newUser := user.NewUser(
		login,
		password,
		name,
		surname,
		patronymic,
		phone,
		email,
	)

	el.mtx.Lock()
	defer el.mtx.Unlock()

	newUser.ID = el.list.GetLen()
	e, err := el.list.AddData(*newUser)
	if err != nil {
		return user.User{}, err
	}

	return e, nil
}

func (ul *UserList) DeleteUser(ctx context.Context, id int) error {
	ul.mtx.Lock()
	defer ul.mtx.Unlock()

	if err := ul.list.DeleteData(id); err != nil {
		if errors.Is(err, ilist.ErrDataNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}
