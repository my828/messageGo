// Test correct inputs
// Test incorrect inputs
// Check response body, status, headers
package handlers

import (
	"assignments-my828/servers/gateway/models/users"
	"assignments-my828/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// catch http errors with header, content type, and code
// func Error(w ResponseWriter, error string, code int) {
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	w.WriteHeader(code)
// 	fmt.Fprintln(w, error)
// }
func TestContext_InsertNewUserHandler(t *testing.T) {
	validUser := users.User{
		ID:        1,
		Email:     "testing@example.com",
		FirstName: "Min",
		LastName:  "Yang",
	}

	validNewUser := users.NewUser{
		Email:        "testing@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "minyang",
		FirstName:    "Min",
		LastName:     "Yang",
	}

	invalidNewUser := users.NewUser{
		Email:        "testing@example.com",
		Password:     "pass",
		PasswordConf: "password",
		UserName:     "minyang",
		FirstName:    "Min",
		LastName:     "Yang",
	}

	cases := []struct {
		name                string
		requestBody         *users.NewUser
		requestType         string
		method              string
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid Post Request",
			&validNewUser,
			"application/json",
			http.MethodPost,
			http.StatusCreated,
			false,
			"application/json",
			&validUser,
		},
		{
			"Invalid Post Request",
			&validNewUser,
			"application/json",
			http.MethodGet,
			http.StatusMethodNotAllowed,
			true,
			"text/plain; charset=utf-8",
			nil,
		},
		{
			"Invalid Content-Type Request",
			&validNewUser,
			"text/plain",
			http.MethodPost,
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			nil,
		},
		{
			"Invalid User Request",
			&invalidNewUser,
			"application/json",
			http.MethodPost,
			500,
			true,
			"text/plain; charset=utf-8",
			nil,
		},
		{
			"Invalid New User Request",
			&users.NewUser{
				FirstName: "fail",
			},
			"application/json",
			http.MethodPost,
			http.StatusBadRequest,
			true,
			"text/plain; charset=utf-8",
			nil,
		},
	}

	for _, c := range cases {
		context := &Context{}
		body, _ := json.Marshal(c.requestBody)
		if c.requestBody.FirstName == "fail" {
			body = []byte("")
		}
		request := httptest.NewRequest(c.method, "/v1/users", bytes.NewBuffer(body))
		request.Header.Set(ContentTypeHeader, c.requestType)
		recorder := httptest.NewRecorder()
		// make sql store
		user, _ := c.requestBody.ToUser()
		context.UsersStore = NewMockStore(user)
		// make mem store
		context.SessionStore = sessions.NewMemStore(time.Hour, time.Minute)
		context.UsersHandler(recorder, request)
		response := recorder.Result()

		resContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != resContentType {
			t.Errorf("case %s: incorrect return type: expected: %s recieved: %s",
				c.name, c.expectedContentType, resContentType)
		}

		resStatusCode := response.StatusCode
		if c.expectedStatusCode != resStatusCode {
			t.Errorf("case %s: incorrect status code: expected: %d recieved: %d",
				c.name, c.expectedStatusCode, resStatusCode)
		}

		user = &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but revieved none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("case %s: incorrect return: expected %v but revieved %v",
				c.name, c.expectedReturn, user)
		}
	}
}

func TestContext_PatchUserHandler(t *testing.T) {
	context := &Context{}
	validUser := users.User{
		ID:        1,
		Email:     "testing@example.com",
		FirstName: "Min",
		LastName:  "Yang",
	}
	update := &users.Updates{
		FirstName: "Success",
		LastName:  "Update",
	}
	cases := []struct {
		name                string
		method              string
		idPath              string
		requestType         string
		requestBody         *users.Updates
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid Patch Request",
			http.MethodPatch,
			"1",
			ContentTypeApplicationJSON,
			update,
			http.StatusOK,
			false,
			"application/json",
			&users.User{
				Email:     "testing@example.com",
				FirstName: "Success",
				LastName:  "Update",
			},
		},
		{
			"Valid Patch Request",
			http.MethodPatch,
			"me",
			ContentTypeApplicationJSON,
			update,
			http.StatusOK,
			false,
			"application/json",
			&users.User{
				Email:     "testing@example.com",
				FirstName: "Success",
				LastName:  "Update",
			},
		},
		{
			"Invalid Specific Request",
			http.MethodPost,
			"me",
			ContentTypeApplicationJSON,
			update,
			http.StatusMethodNotAllowed,
			true,
			"text/plain; charset=utf-8",
			&users.User{
				Email:     "testing@example.com",
				FirstName: "Success",
				LastName:  "Update",
			},
		},
		{
			"Invalid ID Request",
			http.MethodPatch,
			"bad",
			ContentTypeApplicationJSON,
			update,
			http.StatusBadRequest,
			true,
			ContentTypeTextPlain,
			&users.User{},
		},
		{
			"Invalid ID Request",
			http.MethodPatch,
			"2",
			ContentTypeApplicationJSON,
			update,
			http.StatusForbidden,
			true,
			ContentTypeTextPlain,
			&users.User{},
		},
		{
			"Invalid Content-Type Request",
			http.MethodPatch,
			"1",
			ContentTypeTextPlain,
			update,
			http.StatusUnsupportedMediaType,
			true,
			ContentTypeTextPlain,
			&users.User{},
		},
		{
			"Invalid Body Request",
			http.MethodPatch,
			"1",
			ContentTypeApplicationJSON,
			&users.Updates{
				FirstName: "fail",
			},
			http.StatusBadRequest,
			true,
			ContentTypeTextPlain,
			&users.User{},
		},
		{
			"Valid Get Request",
			http.MethodGet,
			"1",
			ContentTypeApplicationJSON,
			update,
			http.StatusOK,
			false,
			"application/json",
			&validUser,
		},
		{
			"Invalid Session Request",
			http.MethodGet,
			"bad",
			ContentTypeApplicationJSON,
			update,
			http.StatusBadRequest,
			true,
			ContentTypeTextPlain,
			&users.User{},
		},
		{
			"Invalid Session Request",
			http.MethodGet,
			"2",
			ContentTypeApplicationJSON,
			update,
			http.StatusNotFound,
			true,
			ContentTypeTextPlain,
			&users.User{},
		},
	}

	for _, c := range cases {
		body, _ := json.Marshal(c.requestBody)
		if c.requestBody.FirstName == "fail" {
			body = []byte("")
		}
		request := httptest.NewRequest(c.method, "/v1/users/"+c.idPath, bytes.NewBuffer(body))
		request.Header.Set(ContentTypeHeader, c.requestType)

		recorder := httptest.NewRecorder()

		context.Key = "key"
		context.SessionStore = sessions.NewMemStore(time.Hour, time.Minute)
		sessionID, _ := sessions.NewSessionID(context.Key)
		request.Header.Add("Authorization", "Bearer "+sessionID.String())

		state := &SessionState{
			time.Now(),
			&validUser,
		}
		context.SessionStore.Save(sessionID, state)
		context.UsersStore = NewMockStore(&validUser)
		context.SpecificUserHandler(recorder, request)
		response := recorder.Result()

		resContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != resContentType {
			t.Errorf("case %s: incorrect return type: expected: %s recieved: %s",
				c.name, c.expectedContentType, resContentType)
		}

		resStatusCode := response.StatusCode
		if c.expectedStatusCode != resStatusCode {
			t.Errorf("case %s: incorrect status code: expected: %d recieved: %d",
				c.name, c.expectedStatusCode, resStatusCode)
		}

		user := &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but revieved none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("case %s: incorrect return: expected %v but revieved %v",
				c.name, c.expectedReturn, user)
		}
	}
}

func TestContext_PostSessionHandler(t *testing.T) {
	context := &Context{}
	validUser := users.NewUser{
		Email:        "testing@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "minyang",
		FirstName:    "Min",
		LastName:     "Yang",
	}
	credentials := &users.Credentials{
		Email:    "testing@example.com",
		Password: "password",
	}
	cases := []struct {
		name                string
		method              string
		requestType         string
		requestBody         *users.Credentials
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			"Valid Post Request",
			http.MethodPost,
			ContentTypeApplicationJSON,
			credentials,
			http.StatusCreated,
			false,
			"application/json",
			&users.User{
				Email:     "testing@example.com",
				FirstName: "Min",
				LastName:  "Yang",
			},
		},
		{
			"Not Post Request",
			http.MethodGet,
			ContentTypeApplicationJSON,
			credentials,
			http.StatusMethodNotAllowed,
			true,
			"text/plain; charset=utf-8",
			&users.User{
				Email:     "testing@example.com",
				FirstName: "Min",
				LastName:  "Yang",
			},
		},
		{
			"Bad Content-Type Request",
			http.MethodPost,
			ContentTypeTextPlain,
			credentials,
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			&users.User{},
		},
		{
			"Bad Content-Type Request",
			http.MethodPost,
			ContentTypeTextPlain,
			&users.Credentials{
				Email: "no",
			},
			http.StatusUnsupportedMediaType,
			true,
			"text/plain; charset=utf-8",
			&users.User{},
		},
		{
			"Bad Credential Request",
			http.MethodPost,
			ContentTypeApplicationJSON,
			&users.Credentials{
				Email: "no",
			},
			http.StatusBadRequest,
			true,
			"text/plain; charset=utf-8",
			&users.User{},
		},
	}

	for _, c := range cases {
		body, _ := json.Marshal(c.requestBody)
		if c.requestBody.Email == "no" {
			body = []byte("")
		}
		request := httptest.NewRequest(c.method, "/v1/users/", bytes.NewBuffer(body))
		request.Header.Set(ContentTypeHeader, c.requestType)

		recorder := httptest.NewRecorder()

		context.Key = "Min"
		context.SessionStore = sessions.NewMemStore(time.Hour, time.Minute)
		user, _ := validUser.ToUser()
		context.UsersStore = NewMockStore(user)
		context.SessionsHandler(recorder, request)
		response := recorder.Result()

		resContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != resContentType {
			t.Errorf("case %s: incorrect return type: expected: %s recieved: %s",
				c.name, c.expectedContentType, resContentType)
		}

		resStatusCode := response.StatusCode
		if c.expectedStatusCode != resStatusCode {
			t.Errorf("case %s: incorrect status code: expected: %d recieved: %d",
				c.name, c.expectedStatusCode, resStatusCode)
		}

		user = &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but revieved none", c.name)
		}

		if !c.expectedError && c.expectedReturn.Email != user.Email &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("case %s: incorrect return: expected %v but revieved %v",
				c.name, c.expectedReturn, user)
		}
	}
}

func TestContext_DeleteSessionHandler(t *testing.T) {
	context := &Context{}
	cases := []struct {
		name                string
		method              string
		idPath              string
		expectedStatusCode  int
		expectedError       bool
		expectedContentType string
	}{
		{
			"Valid Delete Request",
			http.MethodDelete,
			"mine",
			http.StatusOK,
			false,
			ContentTypeTextPlain,
		},
		{
			"Invalid Delete Request",
			http.MethodGet,
			"mine",
			http.StatusMethodNotAllowed,
			true,
			ContentTypeTextPlain,
		},
		{
			"Invalid Delete Request",
			http.MethodDelete,
			"",
			http.StatusForbidden,
			true,
			ContentTypeTextPlain,
		},
	}

	for _, c := range cases {
		request := httptest.NewRequest(c.method, "/v1/users/"+c.idPath, nil)

		request.Header.Set(ContentTypeHeader, ContentTypeApplicationJSON)

		recorder := httptest.NewRecorder()
		sid, _ := sessions.NewSessionID("Min")
		context.Key = "Min"
		context.SessionStore = sessions.NewMemStore(time.Hour, time.Minute)
		context.SessionStore.Save(sid, &SessionState{})
		context.SpecificSessionHandler(recorder, request)
		response := recorder.Result()

		resContentType := response.Header.Get("Content-Type")
		if !c.expectedError && c.expectedContentType != resContentType {
			t.Errorf("case %s: incorrect return type: expected: %s recieved: %s",
				c.name, c.expectedContentType, resContentType)
		}

		resStatusCode := response.StatusCode
		if c.expectedStatusCode != resStatusCode {
			t.Errorf("case %s: incorrect status code: expected: %d recieved: %d",
				c.name, c.expectedStatusCode, resStatusCode)
		}
	}
}
