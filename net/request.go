package net

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HttpRequestParams struct {
	Method  string
	Path    string
	Params  map[string]any
	Headers map[string]string
}

func ExecuteHttpRequest(params HttpRequestParams) (code int, body []byte, err error) {
	jsonStr, err := json.Marshal(params.Params)
	if err != nil {
		return
	}

	request, err := http.NewRequest(params.Method, params.Path, bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}

	for k, v := range params.Headers {
		request.Header.Set(k, v)
	}
	client := http.Client{}
	resp, err := client.Do(request.WithContext(context.TODO()))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	code = resp.StatusCode
	return
}
