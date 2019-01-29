package users

import (
	"crypto/md5"
	"io"
	"log"
	"testing"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
func TestToUser(t *testing.T) {
	cases := []struct {
		name          string
		newUser       *NewUser
		expectedError bool
	}{
		{
			"Invalid email",
			&NewUser{
				"<.min@example.com>",
				"password",
				"password",
				"minyang",
				"Min",
				"Yang",
			},
			true,
		},
		{
			"Invalid password",
			&NewUser{
				"<min@example.com>",
				"pass",
				"pass",
				"minyang",
				"Min",
				"Yang",
			},
			true,
		},
		{
			"Invalid password confirmation",
			&NewUser{
				"<min@example.com>",
				"password",
				"passworD",
				"minyang",
				"Min",
				"Yang",
			},
			true,
		},
		{
			"Invalid user name",
			&NewUser{
				"<min@example.com>",
				"password",
				"password",
				"",
				"Min",
				"Yang",
			},
			true,
		},
		{
			"Invalid user name",
			&NewUser{
				"<min@example.com>",
				"password",
				"password",
				"min yang",
				"Min",
				"Yang",
			},
			true,
		},
		{
			"Check email",
			&NewUser{
				"min <YANG@example.com>",
				"password",
				"password",
				"minyang",
				"Min",
				"Yang",
			},
			false,
		},
	}

	for _, c := range cases {
		//user := &User{}
		user, err := c.newUser.ToUser()
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error: %v", c.name, err)
		}
		//Test photoURL
		h := md5.New()
		if _, err := io.WriteString(h, user.Email); err != nil {
			t.Errorf("error: %v", err)
		}

		if !c.expectedError && user.PhotoURL != gravatarBasePhotoURL+string(h.Sum(nil)) {
			log.Printf("user photo: %s", user.Email)
			log.Printf("Other photo: %s", (gravatarBasePhotoURL + string(h.Sum(nil))))
			t.Errorf("Incorrect photo url")
		}
		if !c.expectedError && user.Email != "yang@example.com" {
			t.Errorf("case %s: email not parsed correctly", c.name)
		}
		user.ApplyUpdates(&Updates{"Tom", "Hanks"})
		if user.FirstName != "Tom" && user.LastName != "Hanks" {
			t.Errorf("Not updating user")
		}
		if err := user.ApplyUpdates(nil); err == nil {
			t.Errorf("Expecting error from update")
		}
	}

}
func TestFullName(t *testing.T) {
	cases := []struct {
		name           string
		user           *User
		expectedOutput string
	}{
		{
			"Check full name",
			&User{
				FirstName: "Min",
				LastName:  "Yang",
			},
			"Min Yang",
		},
		{
			"Check last name",
			&User{
				FirstName: "",
				LastName:  "Yang",
			},
			"Yang",
		},
		{
			"Check empty name",
			&User{
				FirstName: "",
				LastName:  "",
			},
			"",
		},
		{
			"Check first name",
			&User{
				FirstName: "Min",
				LastName:  "",
			},
			"Min",
		},
	}
	for _, c := range cases {
		fullName := c.user.FullName()
		if fullName != c.expectedOutput {
			t.Errorf("case %s: not correct full name", c.name)
		}
	}
}

func TestAuth(t *testing.T) {
	cases := []struct {
		name          string
		user          *User
		input         string
		testPass      string
		expectedError bool
	}{
		{
			"Test incorrect password",
			&User{},
			"password",
			"pass",
			true,
		},
		{
			"Test empty password",
			&User{},
			"password",
			"",
			true,
		},
		{
			"Test correct password",
			&User{},
			"password",
			"password",
			false,
		},
	}
	for _, c := range cases {
		c.user.SetPassword(c.input)
		err := c.user.Authenticate(c.testPass)
		if c.expectedError && err == nil {
			t.Errorf("case %s: not catching error for incorrect password", c.name)
		} else if !c.expectedError && err != nil {
			t.Errorf("case %s: got error for correct password", c.name)
		}
	}
}
