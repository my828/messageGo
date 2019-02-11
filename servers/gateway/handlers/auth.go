package handlers

import (
	"assignments-my828/servers/gateway/models/users"
	"assignments-my828/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

const AuthorizationHeader = "Authorization"
const ContentTypeHeader = "Content-Type"
const ContentTypeApplicationJSON = "application/json"
const ContentTypeTextPlain = "text/plain"
const GETID = "select * from users where id=?"
const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "
const forwardFor = "X-Forwarded-For"

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// insert into SignIn (time, asdf) values (NOW(), asdf)
func NewContext(key string, sessionStore sessions.Store, userStore users.Store) *Context {
	return &Context{
		key,
		sessionStore,
		userStore,
	}
}

func (c *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, fmt.Sprintf("Wrong Content-Type! Must be JSON! %v", http.StatusBadRequest), http.StatusUnsupportedMediaType)
			return
		}
		// get request body
		newUser := &users.NewUser{}
		if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err),
				http.StatusBadRequest)
			return
		}
		// convert new user to user
		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// insert to database
		user, err = c.UsersStore.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		state := &SessionState{
			time.Now(),
			user,
		}

		// being new session for new user
		c.Key = string(user.UserName)
		_, err = sessions.BeginSession(c.Key, c.SessionStore, state, w)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}

	if _, err := sessions.GetState(r, c.Key, c.SessionStore, state); err != nil {
		http.Error(w, fmt.Sprintf("Unable get state: %v", err), http.StatusBadRequest)
		return
	}
	idString := path.Base(r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		idIndex, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Did not provide {id} as a number in /v1/user/{id}, please provide the corrct ID"),
				http.StatusBadRequest)
			return
		}
		user, err := c.UsersStore.GetByID(int64(idIndex))
		if err != nil {
			http.Error(w, fmt.Sprintf("ID provided does not exist in data store"), http.StatusNotFound)
			return
		}

		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		if idString != "me" {
			idIndex, err := strconv.Atoi(idString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Did not provide {id} as a number in /v1/user/{id}, please provide the corrct ID"),
					http.StatusBadRequest)
				return
			}
			if int64(idIndex) != state.User.ID {
				http.Error(w, fmt.Sprint("Incorrect user ID"), http.StatusForbidden)
				return
			}
		}

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, fmt.Sprintf("Wrong Content-Type! Must be JSON! %v", http.StatusBadRequest), http.StatusUnsupportedMediaType)
			return
		}

		update := &users.Updates{}
		if err := json.NewDecoder(r.Body).Decode(update); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err),
				http.StatusBadRequest)
			return
		}

		user, err := c.UsersStore.Update(state.User.ID, update)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating user %v:", err), http.StatusBadRequest)
			return
		}

		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, fmt.Sprintf("Wrong Content-Type! Must be JSON! %v", http.StatusBadRequest), http.StatusUnsupportedMediaType)
			return
		}
		credentials := &users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(credentials); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err),
				http.StatusBadRequest)
			return
		}
		user, err := c.UsersStore.GetByEmail(credentials.Email)
		if err != nil {
			// sleep for 2 seconds
			time.Sleep(time.Second * 2)
			http.Error(w, fmt.Sprintf("Error getting user from data store: %v", err),
				http.StatusUnauthorized)
			return
		}

		if err = user.Authenticate(credentials.Password); err != nil {
			http.Error(w, fmt.Sprintf("Not an authenticated user %v", err),
				http.StatusUnauthorized)
			return
		}

		state := &SessionState{
			time.Now(),
			user,
		}

		ipAddr := r.RemoteAddr
		if r.Header.Get(forwardFor) != "" {
			ipAddr = r.Header.Get(forwardFor)
		}
		//record signin user
		signin := &users.SignIn{
			Id: user.ID,
			Time: time.Now(),
			Ip: ipAddr,
		}
		if err := c.UsersStore.InsertSignin(signin); err != nil {
			http.Error(w, fmt.Sprintf("Error inserting to database: %v", err), 500)
			return
		}
		_, err = sessions.BeginSession(c.Key, c.SessionStore, state, w)
		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		path := path.Base(r.URL.Path)
		if path != "mine" {
			http.Error(w, fmt.Sprintf("Inappropriate resource path"), http.StatusForbidden)
			return
		}
		sessions.EndSession(r, c.Key, c.SessionStore)
		w.Header().Add(ContentTypeHeader, ContentTypeTextPlain)
		w.Write([]byte("Signed out"))
	} else {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}
