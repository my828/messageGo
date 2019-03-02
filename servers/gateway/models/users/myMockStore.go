package users

import (
	"fmt"
)

type MyMockStore struct {
	user *User
	err  bool
}

func NewMockStore(user *User) *MyMockStore {
	return &MyMockStore{
		user,
		false,
	}
}

func (ms *MyMockStore) GetByID(id int64) (*User, error) {
	// We can trigger an error by passing in an id of 2
	if id == -1 {
		fmt.Print("!!!!!!!!")
		return nil, fmt.Errorf("Error getting user with id: %d", id)
	}
	return ms.user, nil
}

func (ms *MyMockStore) GetByEmail(email string) (*User, error) {
	// We can trigger an error by passing in an id of 2
	if email != "testing@example.com" {
		return nil, fmt.Errorf("Error getting user with email: %s", email)
	}
	return ms.user, nil
}
func (ms *MyMockStore) GetByUserName(username string) (*User, error) {
	// We can trigger an error by passing in an id of 2
	if username == "minyang" {
		return nil, fmt.Errorf("Error getting user with name: %s", username)
	}

	return ms.user, nil
}

func (ms *MyMockStore) Insert(user *User) (*User, error) {
	if user.FirstName == "Error" {
		return nil, fmt.Errorf("Error Inserting New User")
	}

	// Assumes that if user without the FirstName field being equal to "Error" will
	// always result in a successful insert
	return ms.user, nil
}

func (ms *MyMockStore) Update(id int64, updates *Updates) (*User, error) {
	if id == 2 {
		return &User{}, fmt.Errorf("error updating database")
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

func (sqls *MyMockStore) InsertSignin(signin *SignIn) error {
	return nil
}
