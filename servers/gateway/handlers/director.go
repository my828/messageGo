package handlers

import (
	"assignments-my828/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
)

type Director func(r *http.Request)

func (c *Context) CustomDirector(targets []*url.URL) Director {
	var counter int32
	counter = 0
	state := &SessionState{}
	mx := sync.RWMutex{}
	mx.Lock()
	defer mx.Unlock()
	return func(r *http.Request) {
		targ := targets[counter%int32(len(targets))]
		r.URL.Host = targ.Host
		r.Host = targ.Host
		r.URL.Scheme = "http"
		atomic.AddInt32(&counter, 1) // note, to be extra safe, we’ll need to use mutexes
		counter++
		_, err := sessions.GetState(r, c.Key, c.SessionStore, state)
		if err != nil {
			r.Header.Del("X-User")
			fmt.Printf("Error getting session state/session unauthorized %v", err)
			return
		}
		// note the modulo (%) operator which maps some integer to range from 0 to
		// len(targets)
		result, err := json.Marshal(state.User)
		if err != nil {
			fmt.Printf("Error encoding session state user %v", err)
			return
		}
		r.Header.Add("X-User", string(result))
		r.Header.Add("X-Forwarded-Host", r.Host)
	}
}
