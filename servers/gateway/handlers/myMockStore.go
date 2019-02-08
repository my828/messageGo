package handlers

import (
	"assignments-my828/servers/gateway/models/users"
	"fmt"
)

type MyMockStore struct {
	user *users.User
	err  bool
}

func NewMockStore(user *users.User) *MyMockStore {
	return &MyMockStore{
		user,
		false,
	}
}

func (ms *MyMockStore) GetByID(id int64) (*users.User, error) {
	// We can trigger an error by passing in an id of 2
	if id == 2 {
		return nil, fmt.Errorf("Error getting user with id: %d", id)
	}

	return ms.user, nil
}

func (ms *MyMockStore) GetByEmail(email string) (*users.User, error) {
	// We can trigger an error by passing in an id of 2
	if email != "testing@example.com" {
		return nil, fmt.Errorf("Error getting user with id: %s", email)
	}
	return ms.user, nil
}
func (ms *MyMockStore) GetByUserName(username string) (*users.User, error) {
	// We can trigger an error by passing in an id of 2
	if username == "minyang" {
		return nil, fmt.Errorf("Error getting user with name: %s", username)
	}

	return ms. user, nil
}

func (ms *MyMockStore) Insert(user *users.User) (*users.User, error) {
	if user.FirstName == "Error" {
		return nil, fmt.Errorf("Error Inserting New User")
	}

	// Assumes that if user without the FirstName field being equal to "Error" will
	// always result in a successful insert
	return ms.user, nil
}

func (ms *MyMockStore) Update(id int64, updates *users.Updates) (*users.User, error) {
	if id == 2 {
		return &users.User{}, fmt.Errorf("error updating database")
	}
	ms.user.FirstName = updates.FirstName
	ms.user.LastName = updates.LastName
	return ms.user, nil
}

func (ms *MyMockStore) Delete(id int64) error {
	if id == 2 {
		return fmt.Errorf("error deleting database")
	}
	return nil
}

func (sqls *MyMockStore) InsertSignin(signin *users.SignIn) error {
	return nil
}
