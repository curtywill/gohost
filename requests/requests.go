package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const base_url = "https://cohost.org/api/v1"

// just get the "data" value out of an http response
// stops us from having to explicity create structs for "result" and "data"
func extractData(body []byte) json.RawMessage {
	var respMap map[string]json.RawMessage
	json.Unmarshal(body, &respMap)

	var resultMap map[string]json.RawMessage
	json.Unmarshal(respMap["result"], &resultMap)

	var data json.RawMessage
	json.Unmarshal(resultMap["data"], &data)

	return data
}

func Fetch(client *http.Client, method, endpoint, cookie string, headers map[string]string, body io.Reader, complex bool) ([]byte, http.Header, error) {
	var res *http.Response
	var data []byte

	url := base_url + endpoint
	method = strings.ToUpper(method)

	if method == http.MethodGet {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, nil, err
		}

		if cookie != "" {
			req.AddCookie(&http.Cookie{Name: "connect.sid", Value: cookie})
		}

		res, err = client.Do(req)
		if err != nil {
			return nil, nil, err
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, nil, err
		}

		defer res.Body.Close()
		// "result" -> "data" nesting is only for trpc endpoints
		if endpoint[1:5] == "trpc" {
			data = extractData(data)
		}

		addToCache(cookie, endpoint, data)
	} else if method == http.MethodPost {
		req, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			return nil, nil, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		if cookie != "" {
			req.AddCookie(&http.Cookie{Name: "connect.sid", Value: cookie})
		}

		res, err = client.Do(req)
		if err != nil {
			return nil, nil, err
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, nil, err
		}

		defer res.Body.Close()
	} else if method == http.MethodPut {
		req, err := http.NewRequest(http.MethodPut, url, body)
		if err != nil {
			return nil, nil, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		req.AddCookie(&http.Cookie{Name: "connect.sid", Value: cookie})

		res, err = client.Do(req)
		if err != nil {
			return nil, nil, err
		}
	}

	if res.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("bad request to endpoint %s: %d", endpoint, res.StatusCode)
	}

	if complex {
		return data, res.Header, nil
	}

	return data, nil, nil
}

func FetchTrpc(client *http.Client, methods, cookie string, headers map[string]string) ([]byte, http.Header, error) {
	methods = fmt.Sprintf("/trpc/%s", methods)
	cachedData := getFromCache(cookie, methods)
	if cachedData == nil {
		return Fetch(client, "GET", methods, cookie, headers, nil, false)
	}
	return cachedData, nil, nil
}
