package handlers

import (
	"assignments-my828/servers/gateway/indexes"
	"assignments-my828/servers/gateway/models/users"
	"assignments-my828/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"
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
func NewContext(key string, sessionStore sessions.Store, userStore users.Store, searchIndex *indexes.Trie) *Context {
	return &Context{
		key,
		sessionStore,
		userStore,
		searchIndex,
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
			http.Error(w, err.Error(), 400)
			return
		}

		// insert to database
		user, err = c.UsersStore.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// add user to trie
		c.SearchIndex.SplitNameAddToTrie(user.UserName, user.ID)
		c.SearchIndex.SplitNameAddToTrie(user.FirstName, user.ID)
		c.SearchIndex.SplitNameAddToTrie(user.LastName, user.ID)

		state := &SessionState{
			time.Now(),
			user,
		}

		// being new session for new user
		_, err = sessions.BeginSession(c.Key, c.SessionStore, state, w)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodGet {
		auth := r.Header.Get(headerAuthorization)
		if auth == "" {
			auth = r.URL.Query().Get(paramAuthorization)
		}

		if !strings.Contains(auth, schemeBearer) {
			http.Error(w, fmt.Sprintf("Unauthorized user!"), http.StatusUnauthorized)
			return
		}
		query := r.FormValue("q")
		if len(query) == 0 {
			http.Error(w, fmt.Sprint("Require query!"), http.StatusBadRequest)
			return
		}
		results := c.SearchIndex.Find(query, 20)
		users := []*users.User{}
		for _, i := range results {
			user, err := c.UsersStore.GetByID(i)
			if err != nil {
				http.Error(w, fmt.Sprint("Error getting user from user store"), http.StatusBadRequest)
				return
			}
			users = append(users, user)
		}
		sort.Slice(users, func(i, j int) bool {
			return users[i].UserName < users[j].UserName
		})
		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (c *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}

	_, err := sessions.GetState(r, c.Key, c.SessionStore, state); 
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable get state: %v", err), http.StatusUnauthorized)
		return
	}
	var userID int64
	id := path.Base(r.URL.Path)
	if id == "me" {
		userID = state.User.ID
	} else {
		userID, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Did not provide {id} as a number in /v1/user/{id}, please provide the corrct ID"),
				http.StatusBadRequest)
			return
		}
	}
	switch r.Method {
	case http.MethodGet:
		user, err := c.UsersStore.GetByID(userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("ID provided does not exist in data store"), http.StatusBadRequest)
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
		if userID != state.User.ID {
			http.Error(w, fmt.Sprint("Incorrect user ID"), http.StatusForbidden)
			return
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

		c.SearchIndex.Remove(state.User.FirstName, state.User.ID)
		c.SearchIndex.Remove(state.User.LastName, state.User.ID)
		c.SearchIndex.Add(update.FirstName, state.User.ID)
		c.SearchIndex.Add(update.LastName, state.User.ID)


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
			Id:   user.ID,
			Time: time.Now(),
			Ip:   ipAddr,
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

// func (c *Context) SearchHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		auth := r.Header.Get(headerAuthorization)
// 		if auth == "" {
// 			auth = r.URL.Query().Get(paramAuthorization)
// 		}
// 		if !strings.Contains(auth, schemeBearer) {
// 			http.Error(w, fmt.Sprintf("Unauthorized user!"), http.StatusUnauthorized)
// 		}
// 		query := r.FormValue("q")
// 		if len(query) == 0 {
// 			http.Error(w, fmt.Sprint("Require query!"), http.StatusBadRequest)
// 		}
// 		results := c.SearchIndex.Find(query, 20)
// 		users := []*users.User{}
// 		for _, i := range results {
// 			user, err := c.UsersStore.GetByID(i)
// 			if err != nil {
// 				http.Error(w, fmt.Sprint("Error getting user from user store"), http.StatusBadRequest)
// 			}
// 			users = append(users, user)
// 		}
// 		sort.Slice(users, func(i, j int) bool {
// 			return users[i].UserName < users[j].UserName
// 		})
// 		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
// 		w.WriteHeader(http.StatusCreated)
// 		if err := json.NewEncoder(w).Encode(users); err != nil {
// 			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
// 		return
// 	}
// }
