package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type matcherFunc func() bool
type RequestParams struct {
	vars     map[string]string
	matchers []matcherFunc
	errors   []string
}

func newRequestParams(request *http.Request) *RequestParams {
	v := &RequestParams{
		vars: mux.Vars(request),
	}

	return v
}

func (v *RequestParams) Require(keys ...string) *RequestParams {
	for _, key := range keys {
		v.matchers = append(v.matchers, func() bool {
			if v, ok := v.vars[key]; ok && v != "" {
				return true
			}
			return false
		})
	}

	return v
}

func (v *RequestParams) addError(key, message string) {
	v.errors = append(v.errors, fmt.Sprintf("%s %s", key, message))
}

func (v *RequestParams) Integer(keys ...string) *RequestParams {
	for _, key := range keys {
		v.matchers = append(v.matchers, func() bool {
			if val, ok := v.vars[key]; ok && val != "" {
				_, err := strconv.Atoi(val)
				if err != nil {
					v.addError(key, "must be integer")
					return false
				}

				return true
			}
			return false
		})
	}

	return v
}

func (v *RequestParams) Validate() bool {
	for _, matcher := range v.matchers {
		if !matcher() {
			return false
		}
	}

	return true
}

func (v *RequestParams) GetValidationErrors() []string {
	return v.errors
}

func (v *RequestParams) GetString(key string) string {
	return v.vars[key]
}

func (v *RequestParams) GetInt(key string) int {
	val, err := strconv.Atoi(v.vars[key])
	if err != nil {
		panic(err)
	}
	return val
}
