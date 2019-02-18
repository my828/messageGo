package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	// Use middle ware
	middleware := NewCorsHandler(handler)
	middleware.ServeHTTP(response, req)
	cases := []struct {
		name            string
		expectedOrigin  string
		expectedMethod  string
		expectedType    string
		expectedControl string
		expectedAge     string
	}{
		{
			"Valid Request",
			"*",
			"GET, PUT, POST, PATCH, DELETE",
			ContentTypeHeader + ", " + AuthorizationHeader,
			AuthorizationHeader,
			"600",
		},
	}

	for _, c := range cases {
		origin := response.HeaderMap.Get(allowOrigin)
		if origin != c.expectedOrigin {
			t.Errorf("Wrong Access Control Allow Origin Handler: expected %s but got %s", origin, "*")
		}

		methods := response.HeaderMap.Get(allowMethod)
		if methods != c.expectedMethod {
			t.Errorf("Wrong Access Control Allow Methods Handler: expected %s but got %s", c.expectedMethod, methods)
		}

		headerType := response.HeaderMap.Get(allowHeader)
		if headerType != c.expectedType {
			t.Errorf("Wrong Access Control Allow Headers Handler: expected %s but got %s", c.expectedType, headerType)
		}

		exposed := response.HeaderMap.Get(exposeHeader)
		if exposed != c.expectedControl {
			t.Errorf("Wrong Access Control Expose Headers Handler: expected %s but got %s", c.expectedControl, exposed)
		}

		maxAge := response.HeaderMap.Get(maxAge)
		if maxAge != c.expectedAge {
			t.Errorf("Wrong Access Control Max Age Handler: expected %s but got %s", c.expectedAge, maxAge)
		}
	}
}

