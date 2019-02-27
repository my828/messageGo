// package handlers

// import (
// 	"assignments-my828/servers/gateway/sessions"
// 	"math/rand"
// 	"encoding/json"
// 	"net/http"
// 	"net/url"
// 	"fmt"
// )

// type Director func(r *http.Request)

// // Custom director
// func (c *Context) CustomDirector(targets []*url.URL) Director {
// 	return func(r *http.Request) {
// 		state := &SessionState{}
// 		targ := targets[rand.Int() % len(targets)]
// 		_, err := sessions.GetState(r, c.Key, c.SessionStore, state)
// 		if err != nil {
// 			r.Header.Del("X-User")
// 			fmt.Sprintf("Error getting session state/session unauthorized %v", err)
// 			return
// 		}
// 		// note the modulo (%) operator which maps some integer to range from 0 to
// 		// len(targets)
// 		j, err := json.Marshal(state.User)
// 		if err != nil {
// 			fmt.Sprintf("Error encoding session state user %v", err)
// 			return
// 		}
// 		r.URL.Host = targ.Host
// 		r.Host = targ.Host
// 		r.URL.Scheme = "http"
// 		r.Header.Add("X-User", string(j))
// 		// if _, err := sessions.GetState(r, c.Key, c.SessionStore, state); err == nil {
// 		// 	user, _ := json.Marshal(state.User)
// 		// 	r.Header.Add("X-User", string(user))
// 		// } else {
// 		// 	r.Header.Del("X-User")
// 		// }
// 		// r.Host = targ.Host
// 		// r.URL.Host = targ.Host
// 		// r.URL.Scheme = "http" // targ.Scheme
// 	}
// }


