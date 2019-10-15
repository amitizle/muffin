package checks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// TODO
// * Use jsonpath notation to verify status based on JSON path in the response (https://github.com/tidwall/gjson)
// * on all config, use better option method, i.e sending an `Option` type to the new struct creation.

// HTTPCheck is a struct that defines the HTTP check.
// It holds the HTTP client (that will be configured with oauth2/simple auth)
// as well as the url (`*url.URL`), HTTP method, payload (body data) and `errorHTTPStatusCodes` (= the
// status codes that will be considered as errornous in the check)
type HTTPCheck struct {
	client               *http.Client
	url                  *url.URL
	method               string
	payload              []byte
	errorHTTPStatusCodes map[int]bool
}

// NewHTTPCheck returns a new Check that implements simple HTTP check.
// It receives an `*http.Client` struct as an argument.
// If the passed `*http.Client` is `nil` then it will use the
// default HTTP client (`*http.DefaultClient`).
func NewHTTPCheck(httpClient *http.Client) *HTTPCheck {
	var checkClient *http.Client
	if httpClient == nil {
		checkClient = http.DefaultClient
	} else {
		checkClient = httpClient
	}
	c := &HTTPCheck{
		client:               checkClient,
		method:               http.MethodGet,
		errorHTTPStatusCodes: map[int]bool{},
	}
	c.useDefaultErrorCodes()
	return c
}

func (check *HTTPCheck) useDefaultErrorCodes() {
	for i := 400; i < 600; i++ {
		if http.StatusText(i) != "" {
			check.errorHTTPStatusCodes[i] = true
		}
	}
}

// Run runs the HTTP check
func (check *HTTPCheck) Run() ([]byte, error) {
	req, err := http.NewRequest(check.method, check.url.String(), bytes.NewBuffer(check.payload))
	if err != nil {
		return []byte{}, err
	}

	resp, err := check.client.Do(req)
	if err != nil && resp == nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if check.errorHTTPStatusCodes[resp.StatusCode] {
		return body, fmt.Errorf(resp.Status)
	}

	return body, nil
}

// SetErrorStatusCodes sets the HTTP status codes that'll fail the check.
// For example, if calling `SetErrorStatusCodes([]int{500, 503})` will make the
// check only return `error` (= fail) when getting one of those HTTP statuses back.
func (check *HTTPCheck) SetErrorStatusCodes(codes []int) error {
	badHTTPCodes := []int{}
	for _, statusCode := range codes {
		if statusCode < 100 || statusCode > 599 {
			badHTTPCodes = append(badHTTPCodes, statusCode)
			continue
		}
		check.errorHTTPStatusCodes[statusCode] = true
	}
	if len(badHTTPCodes) > 0 {
		return fmt.Errorf("bad http codes: %v", badHTTPCodes)
	}
	return nil
}

// SetURL sets the URL to which the HTTP check will make the HTTP request.
// It receives a `*url.URL` struct, it's the responsibility of the user of
// this check to build this struct.
func (check *HTTPCheck) SetURL(u *url.URL) error {
	check.url = u
	return nil
}

// GetFullURL returns the string represantation of the check's URL
func (check *HTTPCheck) GetFullURL() string {
	return check.url.String()
}
