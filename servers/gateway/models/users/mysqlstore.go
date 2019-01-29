package users

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SQLStore struct {
	Db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db,
	}
}

//GetByID returns the User with the given ID
func (sqls *SQLStore) GetByID(id int64) (*User, error) {
	rows, err := sqls.Db.Query("select * from users where id=?", id)
	if err != nil {
		return &User{}, fmt.Errorf("Error getting data by id: %v", err)
	}
	defer rows.Close()
	user := &User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return &User{}, fmt.Errorf("Error getting data")
		}
	}
	if err := rows.Err(); err != nil {
		return &User{}, fmt.Errorf("error getting next row: %v\n", err)
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (sqls *SQLStore) GetByEmail(email string) (*User, error) {
	rows, err := sqls.Db.Query("select * from users where email=?", email)
	if err != nil {
		return &User{}, fmt.Errorf("Error getting data by email")
	}
	defer rows.Close()
	user := &User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return &User{}, fmt.Errorf("Error getting data")
		}
	}
	if err := rows.Err(); err != nil {
		return &User{}, fmt.Errorf("error getting next row: %v\n", err)
	}
	return user, nil
}

//GetByUserName returns the User with the given Username
func (sqls *SQLStore) GetByUserName(username string) (*User, error) {
	rows, err := sqls.Db.Query("select * from users where user_name=?", username)
	if err != nil {
		return &User{}, fmt.Errorf("Error getting data by user name")
	}

	defer rows.Close()
	user := &User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return &User{}, fmt.Errorf("Error getting data")
		}
	}
	if err := rows.Err(); err != nil {
		return &User{}, fmt.Errorf("error getting next row: %v\n", err)
	}
	return user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (sqls *SQLStore) Insert(user *User) (*User, error) {
	// insq := "insert into users(id, email, pass_hash, user_name, first_name, last_name
	// 		photo_url) values (?,?,?,?,?,?,?)"
	// stmt, err := sqls.Db.Prepare("insert into users(id, email, pass_hash, user_name, first_name, last_name, photo_url) values (?,?,?,?,?,?,?)")
	// if err != nil {
	// 	return &User{}, fmt.Errorf("error preparing statement %v", err)
	// }
	insq := "insert into users(id, email, pass_hash, user_name, first_name, last_name, photo_url) values (?,?,?,?,?,?,?)"
	res, err := sqls.Db.Exec(insq, user.ID, user.Email, user.PassHash, user.UserName,
		user.FirstName, user.LastName, user.PhotoURL)
	u := &User{
		user.ID,
		user.Email,
		user.PassHash,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.PhotoURL,
	}
	if err != nil {
		return &User{}, fmt.Errorf("error inserting new row: %v\n", err)
	}
	//get the auto-assigned ID for the new row
	id, err := res.LastInsertId()

	u.ID = id
	if err != nil {
		return &User{}, fmt.Errorf("error getting new ID: %v\n", id)
	}
	return u, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (sqls *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "update users set first_name=?, last_name=? where id=?"
	_, err := sqls.Db.Exec(insq, updates.FirstName, updates.LastName, id)
	if err != nil {
		return &User{}, fmt.Errorf("error updating database %v", err)
	}
	return sqls.GetByID(id)
}

//Delete deletes the user with the given ID
func (sqls *SQLStore) Delete(id int64) error {
	stmt, err := sqls.Db.Prepare("delete from users where id=?")
	if err != nil {
		return fmt.Errorf("error preparing database %v", err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error deleteing from database %v", err)
	}
	return nil
}
