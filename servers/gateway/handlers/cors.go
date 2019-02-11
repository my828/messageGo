package handlers

import "net/http"

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/
const allowOrigin = "Access-Control-Allow-Origin"
const allowMethod = "Access-Control-Allow-Methods"
const allowHeader = "Access-Control-Allow-Headers"
const exposeHeader = "Access-Control-Expose-Headers"
const maxAge = "Access-Control-Max-Age"

type Cors struct {
	MyHandler http.Handler
}

func NewCorsHandler(handlerToWrap http.Handler) *Cors {
	return &Cors{handlerToWrap}
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(allowOrigin, "*")
	w.Header().Set(allowMethod, "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set(allowHeader, "Content-Type, Authorization")
	w.Header().Set(exposeHeader, "Authorization")
	w.Header().Set(maxAge, "600")
	c.MyHandler.ServeHTTP(w, r)
}
