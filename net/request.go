package net

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ContentType constants
const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type HttpRequestParams struct {
	Method      string
	Path        string
	Params      map[string]any
	Headers     map[string]string
	ContentType string
}

func ExecuteHttpRequest(params HttpRequestParams) (code int, body []byte, err error) {
	var reqBody *bytes.Buffer

	// 如果未指定ContentType，默认使用JSON
	if params.ContentType == "" {
		params.ContentType = ContentTypeJSON
	}

	// 根据内容类型处理请求参数
	if params.Params != nil {
		if params.ContentType == ContentTypeJSON {
			// JSON格式
			jsonData, err := json.Marshal(params.Params)
			if err != nil {
				return 0, nil, err
			}
			reqBody = bytes.NewBuffer(jsonData)
		} else if params.ContentType == ContentTypeFormURLEncoded {
			// x-www-form-urlencoded格式
			form := url.Values{}
			for key, val := range params.Params {
				// 将任意类型转换为字符串
				strVal, ok := val.(string)
				if ok {
					form.Add(key, strVal)
				} else {
					// 对于非字符串类型，尝试JSON编码
					jsonVal, err := json.Marshal(val)
					if err != nil {
						return 0, nil, err
					}
					form.Add(key, string(jsonVal))
				}
			}
			reqBody = bytes.NewBufferString(form.Encode())
		} else {
			// 不支持的内容类型
			return 0, nil, http.ErrNotSupported
		}
	} else {
		reqBody = &bytes.Buffer{}
	}

	request, err := http.NewRequest(params.Method, params.Path, reqBody)
	if err != nil {
		return
	}

	// 设置Content-Type请求头
	request.Header.Set("Content-Type", params.ContentType)

	if params.Headers != nil {
		for k, v := range params.Headers {
			request.Header.Set(k, v)
		}
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
