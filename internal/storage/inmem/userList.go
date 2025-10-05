package inmem

import (
	"context"
	"errors"
	"practice-backend/internal/models/user"
	"practice-backend/internal/storage/inmem/ilist"
	"sync"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
)

// Concurrent-Use
type UserList struct {
	list        ilist.List[user.User]
	loginToUser map[string]*user.User
	mtx         *sync.Mutex
}

func NewUserList() UserList {
	return UserList{
		list:        ilist.NewList[user.User](),
		loginToUser: make(map[string]*user.User),
		mtx:         new(sync.Mutex),
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

func (ul *UserList) GetUserByLogin(ctx context.Context, login string) (user.User, error) {
	ul.mtx.Lock()
	defer ul.mtx.Unlock()

	usr, ok := ul.loginToUser[login]
	if !ok {
		return user.User{}, ErrUserNotFound
	}
	return *usr, nil
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
	isAdmin bool,
) (user.User, error) {
	newUser := user.NewUser(
		login,
		password,
		name,
		surname,
		patronymic,
		phone,
		email,
		isAdmin,
	)

	if _, err := el.GetUserByLogin(ctx, login); err == nil {
		return user.User{}, ErrUserAlreadyExist
	}

	el.mtx.Lock()
	defer el.mtx.Unlock()

	newUser.ID = el.list.GetLen()

	e, err := el.list.AddData(*newUser)
	if err != nil {
		return user.User{}, err
	}
	el.loginToUser[newUser.Login] = newUser

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
