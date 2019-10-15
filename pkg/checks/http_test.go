package checks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getServer(respStatusCode int, respBody string) *httptest.Server {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if respStatusCode >= 400 {
			http.Error(w, http.StatusText(respStatusCode), respStatusCode)
			return
		}
		fmt.Fprintln(w, respBody)
	}))
	return ts
}

func TestHTTPRun(t *testing.T) {
	tests := []struct {
		name                  string
		serverStatusCode      int
		serverResponse        string
		checkErrorStatusCodes []int
		shouldFailCheck       bool
	}{
		{
			name:                  "test that 500 return error when it's configured with 500 as an error",
			serverStatusCode:      http.StatusInternalServerError,
			checkErrorStatusCodes: []int{http.StatusInternalServerError},
			shouldFailCheck:       true,
		},
		{
			name:                  "test that 500 return success when it's configured without 500 as an error",
			serverStatusCode:      http.StatusInternalServerError,
			checkErrorStatusCodes: []int{},
			shouldFailCheck:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := getServer(test.serverStatusCode, test.serverResponse)
			defer srv.Close()

			httpCheck := NewHTTPCheck(srv.Client())
			u, _ := url.Parse(srv.URL)
			httpCheck.SetURL(u)
			httpCheck.SetErrorStatusCodes(test.checkErrorStatusCodes)
			_, err := httpCheck.Run()
			if test.shouldFailCheck && err == nil {
				t.Errorf("test should have failed but succeeded: %v", err)
			}
		})
	}
}

func TestDefaultHTTPErrorCodes(t *testing.T) {
	c := NewHTTPCheck(nil)
	if len(c.errorHTTPStatusCodes) == 0 {
		t.Errorf("expected errorHTTPStatusCodes to have more than zero default HTTP status error codes")
	}

	if c.errorHTTPStatusCodes[599] {
		t.Error("expected HTTP status 599 to not be included in the default HTTP status error codes")
	}

	if !c.errorHTTPStatusCodes[http.StatusNotFound] {
		t.Errorf("expected HTTP status %d to be included in the default HTTP status error codes", http.StatusNotFound)
	}
}

func TestSetURL(t *testing.T) {
}
