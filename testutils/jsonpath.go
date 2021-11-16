package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/PaesslerAG/jsonpath"
	"github.com/steinfletcher/apitest"
)

func Extract(expression string, value interface{}) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		vr := reflect.ValueOf(value)
		if vr.Kind() != reflect.Ptr || vr.IsNil() {
			return fmt.Errorf("value isn't a pointer")
		}

		r, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		vr.Elem().Set(reflect.ValueOf(r))
		return nil
	}
}

// from https://github.com/steinfletcher/apitest-jsonpath/blob/master/jsonpath.go
func jsonPath(reader io.Reader, expression string) (interface{}, error) {
	v := interface{}(nil)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(expression, v)
	if err != nil {
		return nil, fmt.Errorf("evaluating '%s' resulted in error: '%s'", expression, err)
	}
	return value, nil
}
