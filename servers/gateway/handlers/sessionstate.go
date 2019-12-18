package handlers

import (
	"assignments-my828/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

type SessionState struct {
	SessionBegin time.Time   `json:"sessionBegin,omitempty"`
	User         *users.User `json:"user,omitempty"`
}
