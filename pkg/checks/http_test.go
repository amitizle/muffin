package checks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testCtx = context.Background()
)

func getServer(respStatusCode int, respBody string) *httptest.Server {
	// httptest.NewServer
	// ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		configInput           map[string]interface{}
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
			checkErrorStatusCodes: []int{400},
			shouldFailCheck:       false,
		},
		{
			name:                  "non existing method - should fail running",
			serverStatusCode:      http.StatusInternalServerError,
			checkErrorStatusCodes: []int{400},
			shouldFailCheck:       true,
			configInput:           map[string]interface{}{"method": "///"},
		},
		{
			name:                  "non existing URL - should fail making the request",
			serverStatusCode:      http.StatusInternalServerError,
			checkErrorStatusCodes: []int{400},
			shouldFailCheck:       true,
			configInput:           map[string]interface{}{"url": "http://www.abcdcsdkfjhskfjhkjrewhtkjdhfkjdhfgkjehr.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := getServer(test.serverStatusCode, test.serverResponse)
			defer srv.Close()

			httpCheck := &HTTPCheck{}
			httpConfigMap := test.configInput
			if httpConfigMap == nil {
				httpConfigMap = map[string]interface{}{}
			}
			_, ok := httpConfigMap["url"]
			if !ok {
				httpConfigMap["url"] = srv.URL
			}
			httpConfigMap["error_http_status_codes"] = test.checkErrorStatusCodes

			httpCheck.Initialize(testCtx)
			httpCheck.Configure(httpConfigMap)
			_, err := httpCheck.Run()
			if test.shouldFailCheck && err == nil {
				t.Fatalf("test should have failed but succeeded")
			}

			if !test.shouldFailCheck && err != nil {
				t.Fatalf("failed running HTTP check: %v", err)
			}
		})
	}
}

func TestDefaultHTTPErrorCodes(t *testing.T) {
	c := &HTTPCheck{}
	c.Initialize(testCtx)
	c.Configure(map[string]interface{}{"url": "http://www.example.com"})
	conf := c.config
	if len(conf.errorHTTPStatusCodesMap) == 0 {
		t.Errorf("expected errorHTTPStatusCodes to have more than zero default HTTP status error codes")
	}

	if conf.errorHTTPStatusCodesMap[599] {
		t.Error("expected HTTP status 599 to not be included in the default HTTP status error codes")
	}

	if !conf.errorHTTPStatusCodesMap[http.StatusNotFound] {
		t.Errorf("expected HTTP status %d to be included in the default HTTP status error codes", http.StatusNotFound)
	}
}

func TestConfigure(t *testing.T) {
	tests := []struct {
		name                string
		expectedConfig      *HTTPCheckConfig
		configInput         map[string]interface{}
		shouldFailConfigure bool
	}{
		{
			name: "no URL given - configuration should fail",
			expectedConfig: &HTTPCheckConfig{
				Method:               http.MethodGet,
				ErrorHTTPStatusCodes: []int{400},
			},
			configInput: map[string]interface{}{
				"method":                  "GET",
				"error_http_status_codes": []int{400},
			},
			shouldFailConfigure: true,
		},
		{
			name: "all configuration given, all legal",
			expectedConfig: &HTTPCheckConfig{
				URL:                  "http://www.example.com",
				Method:               http.MethodGet,
				ErrorHTTPStatusCodes: []int{400},
			},
			configInput: map[string]interface{}{
				"url":                     "http://www.example.com",
				"method":                  "GET",
				"error_http_status_codes": []int{400},
			},
			shouldFailConfigure: false,
		},
		{
			name: "no HTTP method given - should use the default method",
			expectedConfig: &HTTPCheckConfig{
				URL:                  "http://www.example.com",
				Method:               http.MethodHead,
				ErrorHTTPStatusCodes: []int{400},
			},
			configInput: map[string]interface{}{
				"url":                     "http://www.example.com",
				"error_http_status_codes": []int{400},
			},
			shouldFailConfigure: false,
		},
		{
			name: "HTTP error status codes not given - should use default ones",
			expectedConfig: &HTTPCheckConfig{
				URL:                  "http://www.example.com",
				Method:               http.MethodGet,
				ErrorHTTPStatusCodes: []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 419, 420, 421, 422, 423, 424, 425, 426, 427, 428, 429, 430, 431, 432, 433, 434, 435, 436, 437, 438, 439, 440, 441, 442, 443, 444, 445, 446, 447, 448, 449, 450, 451, 452, 453, 454, 455, 456, 457, 458, 459, 460, 461, 462, 463, 464, 465, 466, 467, 468, 469, 470, 471, 472, 473, 474, 475, 476, 477, 478, 479, 480, 481, 482, 483, 484, 485, 486, 487, 488, 489, 490, 491, 492, 493, 494, 495, 496, 497, 498, 499, 500, 501, 502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 512, 513, 514, 515, 516, 517, 518, 519, 520, 521, 522, 523, 524, 525, 526, 527, 528, 529, 530, 531, 532, 533, 534, 535, 536, 537, 538, 539, 540, 541, 542, 543, 544, 545, 546, 547, 548, 549, 550, 551, 552, 553, 554, 555, 556, 557, 558, 559, 560, 561, 562, 563, 564, 565, 566, 567, 568, 569, 570, 571, 572, 573, 574, 575, 576, 577, 578, 579, 580, 581, 582, 583, 584, 585, 586, 587, 588, 589, 590, 591, 592, 593, 594, 595, 596, 597, 598, 599},
			},
			configInput: map[string]interface{}{
				"url":    "http://www.example.com",
				"method": "GET",
			},
			shouldFailConfigure: false,
		},
		{
			name: "incorrect types in config map - should fail configure",
			expectedConfig: &HTTPCheckConfig{
				URL:                  "http://www.example.com",
				Method:               http.MethodHead,
				ErrorHTTPStatusCodes: []int{400},
			},
			configInput: map[string]interface{}{
				"url":    "http://www.example.com",
				"method": 1234,
			},
			shouldFailConfigure: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &HTTPCheck{}
			c.Initialize(testCtx)
			err := c.Configure(test.configInput)
			if err != nil && !test.shouldFailConfigure {
				t.Fatalf("expected configuration to succeed, configuration failed: %v", err)
			} else if err != nil && test.shouldFailConfigure {
				return
			}
			parsedConf := c.config
			if parsedConf.parsedURL.String() != test.expectedConfig.URL {
				t.Fatalf("expected parsed url (%s) to be equal to the expected config (%s)", parsedConf.parsedURL.String(), test.expectedConfig.URL)
			}
		})
	}
}

func TestGetFullURL(t *testing.T) {
	url := "http://www.example.com"
	c := &HTTPCheck{}
	c.Initialize(testCtx)
	c.Configure(map[string]interface{}{"url": url})

	if c.GetFullURL() != url {
		t.Fatalf("expected to get full url %s, got %s", url, c.GetFullURL())
	}
}
