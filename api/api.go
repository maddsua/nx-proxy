package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type APIResponse[T any] struct {
	Data  *T        `json:"data"`
	Error *APIError `json:"error"`
}

func DecodeAPIResponse[T any](reader io.Reader) (*APIResponse[T], error) {

	var val APIResponse[T]
	if err := json.NewDecoder(reader).Decode(&val); err != nil && err != io.EOF {
		return nil, err
	}

	return &val, nil
}

type APIError struct {
	Message string `json:"message"`
}

func beacon(client *Client, method string, path string, payload any) error {
	if _, err := fetch[any](client, method, path, payload); err != nil {
		return err
	}
	return nil
}

func fetch[T any](client *Client, method string, path string, payload any) (*T, error) {

	if client.URL == nil {
		return nil, fmt.Errorf("remote url not set")
	}

	reqUrl := url.URL{
		Scheme:   client.URL.Scheme,
		Host:     client.URL.Host,
		Path:     strings.TrimRight(client.URL.Path, "/") + path,
		RawQuery: client.URL.RawQuery,
	}

	var bodyReader io.Reader
	if payload != nil {
		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(path); err != nil {
			return nil, fmt.Errorf("marshal: %v", err)
		}
		bodyReader = &buff
	}

	req, err := http.NewRequest(method, reqUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	if client.Token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token.String()))
	}

	req = nil

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil, nil
	}

	apiResp, err := DecodeAPIResponse[T](resp.Body)
	if err != nil {
		return nil, err
	}

	if apiResp.Error != nil {
		return nil, errors.New(apiResp.Error.Message)
	} else if apiResp.Data == nil && (resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusResetContent) {
		return nil, fmt.Errorf("http: %s", resp.Status)
	}

	return apiResp.Data, nil
}
