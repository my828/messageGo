package users

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const GETID = "select * from users where id=?"
const GETEMAIL = "select * from users where email=?"
const GETUSER = "select * from users where user_name=?"
const INSERTSTATEMENT = "insert into users(id, email, pass_hash, user_name, first_name, last_name, photo_url) values (?,?,?,?,?,?,?)"
const UPDATESTATEMENT = "update users set first_name=?, last_name=? where id=?"
const DELETESTATEMENT = "delete from users where id=?"

func TestMySQLStore_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	expectedUser := &User{
		FirstName: "Min",
		LastName:  "Yang",
	}
	store := NewSQLStore(db)

	// Create a row with the appropriate fields in your SQL database
	// Add the actual values to the row
	row := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	// Expecting a successful "query"
	// This tells our db to expect this query (id) as well as supply a certain response (row)
	// REMINDER: Since sqlmock requires a regex string, in order for `?` to be interpreted, you'll
	// have to wrap it within a `regexp.QuoteMeta`. Be mindful that you will need to do this EVERY TIME you're
	// using any reserved metacharacters in regex.
	mock.ExpectQuery(regexp.QuoteMeta(GETID)).
		WithArgs(expectedUser.ID).WillReturnRows(row)

	// Since we know our query is successful, we want to test whether there happens to be
	// any expected error that may occur.
	user, err := store.GetByID(expectedUser.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Again, since we are assuming that our query is successful, we can test for when our
	// function doesn't work as expected.
	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User queried does not match expected user")
	}

	// Expecting a unsuccessful "query"
	// Attempting to search by an id that doesn't exist. This would result in a
	// sql.ErrNoRows error
	// REMINDER: Using a constant makes your code much clear, and is highly recommended.
	mock.ExpectQuery(regexp.QuoteMeta(GETID)).
		WithArgs(-1).WillReturnError(sql.ErrNoRows)

	// Since we are expecting an error here, we create a condition opposing that to see
	// if our GetById is working as expected
	if _, err = store.GetByID(-1); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	// Attempting to trigger a DBMS querying error
	queryingErr := fmt.Errorf("DBMS error when querying")

	if _, err = store.GetByID(expectedUser.ID); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", queryingErr)
	}

	wrongRow := sqlmock.NewRows([]string{"id", "pass_hash", "first_name", "last_name", "photo_url"})
	wrongRow.AddRow(expectedUser.ID, expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	mock.ExpectQuery(regexp.QuoteMeta(GETID)).
		WithArgs(expectedUser.ID).WillReturnRows(wrongRow)
	_, err = store.GetByID(expectedUser.ID)
	if err == nil {
		t.Errorf("Expected error: %v but recieved nil", queryingErr)
	}

	// This attempts to check if there are any expectations that we haven't met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}

func TestMySQLStore_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	expectedUser := &User{
		FirstName: "Min",
		LastName:  "Yang",
	}
	store := NewSQLStore(db)

	row := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(GETEMAIL)).
		WithArgs(expectedUser.Email).WillReturnRows(row)

	user, err := store.GetByEmail(expectedUser.Email)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User queried does not match expected user")
	}

	mock.ExpectQuery(regexp.QuoteMeta(GETEMAIL)).
		WithArgs("").WillReturnError(sql.ErrNoRows)

	if _, err = store.GetByEmail(""); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	queryingErr := fmt.Errorf("DBMS error when querying")

	wrongRow := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "photo_url"})
	wrongRow.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	mock.ExpectQuery(regexp.QuoteMeta(GETEMAIL)).
		WithArgs(expectedUser.Email).WillReturnRows(wrongRow)
	_, err = store.GetByEmail(expectedUser.Email)
	if err == nil {
		t.Errorf("Expected error: %v but recieved nil", queryingErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}
func TestMySQLStore_GetByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	expectedUser := &User{
		FirstName: "Min",
		LastName:  "Yang",
	}
	store := NewSQLStore(db)

	row := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(GETUSER)).
		WithArgs(expectedUser.UserName).WillReturnRows(row)

	user, err := store.GetByUserName(expectedUser.UserName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User queried does not match expected user")
	}

	mock.ExpectQuery(regexp.QuoteMeta(GETUSER)).
		WithArgs("").WillReturnError(sql.ErrNoRows)

	if _, err = store.GetByUserName(""); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	queryingErr := fmt.Errorf("DBMS error when querying")

	wrongRow := sqlmock.NewRows([]string{"id", "email", "user_name", "first_name", "last_name", "photo_url"})
	wrongRow.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	mock.ExpectQuery(regexp.QuoteMeta(GETUSER)).
		WithArgs(expectedUser.UserName).WillReturnRows(wrongRow)
	_, err = store.GetByUserName(expectedUser.UserName)
	if err == nil {
		t.Errorf("Expected error: %v but recieved nil", queryingErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}

func TestMySQLStore_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	inputUser := &User{
		FirstName: "Min",
		LastName:  "Yang",
	}
	store := NewSQLStore(db)
	// This tells our db to expect an insert query with certain arguments with a certain
	// return result
	mock.ExpectExec(regexp.QuoteMeta("insert into users(id, email, pass_hash, user_name, first_name, last_name, photo_url) values (?,?,?,?,?,?,?)")).
		WithArgs(inputUser.ID, inputUser.Email, inputUser.PassHash, inputUser.UserName, inputUser.FirstName, inputUser.LastName, inputUser.PhotoURL).
		WillReturnResult(sqlmock.NewResult(2, 1))

	user, err := store.Insert(inputUser)
	inputUser.ID = user.ID
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(user, inputUser) {
		t.Errorf("User returned does not match input user")
	}

	invalidUser := &User{
		-1,
		"min@example.com",
		nil,
		"minyang",
		"min",
		"Yang",
		"hello",
	}

	insertErr := fmt.Errorf("Error executing INSERT operation")

	if _, err = store.Insert(invalidUser); err == nil {
		t.Errorf("Expected error: %v but recieved nil", insertErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}
func TestMySQLStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	expectedUser := &User{
		FirstName: "Min",
		LastName:  "Yang",
	}
	store := NewSQLStore(db)

	update := &Updates{
		"Tom",
		"Hanks",
	}

	row := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectExec(regexp.QuoteMeta(UPDATESTATEMENT)).
		WithArgs(update.FirstName, update.LastName, expectedUser.ID).
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectQuery(regexp.QuoteMeta("select * from users where id=?")).
		WithArgs(expectedUser.ID).
		WillReturnRows(row)

	updatedUser, err := store.Update(expectedUser.ID, update)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(expectedUser, updatedUser) {
		t.Errorf("User returned does not match input user")
	}

	noUpdate := &Updates{
		"Min",
		"Yang",
	}

	updateErr := fmt.Errorf("Error executing UPDATE operation")
	mock.ExpectExec(regexp.QuoteMeta(UPDATESTATEMENT)).
		WithArgs(noUpdate.FirstName, noUpdate.LastName, updatedUser.ID).
		WillReturnResult(sqlmock.NewResult(2, 1))
	if _, err = store.Update(updatedUser.ID, noUpdate); err == nil {
		t.Errorf("Expected error: %v but recieved nil", updateErr)
	}
	invalidUpdate := &Updates{}

	if _, err = store.Update(updatedUser.ID, invalidUpdate); err == nil {
		t.Errorf("Expected error: %v but recieved nil", updateErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}
func TestMySQLStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	expectedUser := &User{
		FirstName: "Min",
		LastName:  "Yang",
	}
	row := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	store := NewSQLStore(db)

	mock.ExpectExec(regexp.QuoteMeta(DELETESTATEMENT)).
		WithArgs(expectedUser.ID).
		WillReturnResult(sqlmock.NewResult(2, 1))
	err = store.Delete(expectedUser.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	deleteErr := fmt.Errorf("Error executing DELETE operation")

	if err = store.Delete(-1); err == nil {
		t.Errorf("Expected error: %v but recieved nil", deleteErr)
	}
}
