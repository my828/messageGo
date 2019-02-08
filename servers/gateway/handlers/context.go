package handlers

import (
	"assignments-my828/servers/gateway/models/users"
	"assignments-my828/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

type Context struct {
	Key          string
	SessionStore sessions.Store
	UsersStore   users.Store
}
