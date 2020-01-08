package checks

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/amitizle/muffin/internal/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
)

// TODO
// * Use jsonpath notation to verify status based on JSON path in the response (https://github.com/tidwall/gjson)
// * on all config, use better option method, i.e sending an `Option` type to the new struct creation.
// * Support configuring headers (such as auth headers)

// HTTPCheck is a struct that defines the HTTP check.
// It holds the HTTP client and an HTTPCheckConfig struct.
type HTTPCheck struct {
	client *http.Client
	config *HTTPCheckConfig
	ctx    context.Context
	logger zerolog.Logger
}

// HTTPCheckConfig is a struct that holds the configuration required for the HTTP check.
// It is populated with `mapstructure` and holds some private fields that suppose to
// hold a parsed/verified version of the configuration input.
type HTTPCheckConfig struct {
	URL                  string `mapstructure:"url"`
	Method               string `mapstructure:"method"`
	Payload              []byte `mapstructure:"payload"`
	ErrorHTTPStatusCodes []int  `mapstructure:"error_http_status_codes"`

	// private fields
	errorHTTPStatusCodesMap map[int]bool
	parsedURL               *url.URL
}

// useDefaultErrorCodes populated an HTTPCheckConfig's HTTP error codes
// with default ones (400 - 599)
func (checkConfig *HTTPCheckConfig) useDefaultErrorCodes() {
	for i := 400; i < 600; i++ {
		if http.StatusText(i) != "" {
			checkConfig.errorHTTPStatusCodesMap[i] = true
		}
	}
}

// Initialize initializing an HTTP client for the HTTPCheck.
func (check *HTTPCheck) Initialize(ctx context.Context) error {
	check.client = http.DefaultClient
	check.ctx = ctx
	lg, err := logger.GetContext(ctx)
	if err != nil {
		return err
	}
	check.logger = lg
	return nil
}

// Configure decodes map[string]interface{} to an HTTPCheckConfig struct instance.
// It does so using `mapstructure`.
// After decoding, it configures some default values in case they were not given in
// the configuration.
func (check *HTTPCheck) Configure(config map[string]interface{}) error {
	httpConfig := &HTTPCheckConfig{
		errorHTTPStatusCodesMap: map[int]bool{},
	}
	if err := mapstructure.Decode(config, httpConfig); err != nil {
		return err
	}
	u, err := url.ParseRequestURI(httpConfig.URL)
	if err != nil {
		return err
	}
	httpConfig.parsedURL = u

	if httpConfig.ErrorHTTPStatusCodes != nil {
		for _, errStatusCode := range httpConfig.ErrorHTTPStatusCodes {
			httpConfig.errorHTTPStatusCodesMap[errStatusCode] = true
		}
	} else {
		httpConfig.useDefaultErrorCodes()
	}

	if httpConfig.Method == "" {
		httpConfig.Method = http.MethodHead
	}

	if httpConfig.Payload == nil {
		httpConfig.Payload = []byte{}
	}

	check.config = httpConfig
	return nil
}

// Run runs the HTTP check
func (check *HTTPCheck) Run() ([]byte, error) {
	check.logger.Debug().Msg("running check")
	req, err := http.NewRequest(check.config.Method, check.config.parsedURL.String(), bytes.NewBuffer(check.config.Payload))
	if err != nil {
		check.logger.Error().Err(err).Msg("check encountered an error")
		return []byte{}, err
	}

	resp, err := check.client.Do(req)
	// if resp is not nil it means that the HTTP request failed, however the
	// check itself should be reporting an error, thus we won't return nil
	if err != nil && resp == nil {
		check.logger.Error().Err(err).Msg("check encountered an error")
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.logger.Error().Err(err).Msg("check encountered an error")
		return []byte{}, err
	}

	if check.config.errorHTTPStatusCodesMap[resp.StatusCode] {
		return body, fmt.Errorf(resp.Status)
	}

	return body, nil
}

// GetFullURL returns the string represantation of the check's URL
func (check *HTTPCheck) GetFullURL() string {
	return check.config.parsedURL.String()
}
