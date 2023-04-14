package requests

import (
	"encoding/json"
	"fmt"
	"gohost/structs"
	"io"
	"log"
	"net/http"
	"strings"
)

const base_url = "https://cohost.org/api/v1"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// just get the "data" value out of an http response
// stops us from having to explicity create structs for "result" and "data"
func extractData(body []byte) json.RawMessage {
	var respMap map[string]json.RawMessage
	err := json.Unmarshal(body, &respMap)
	check(err)

	var resultMap map[string]json.RawMessage
	err = json.Unmarshal(respMap["result"], &resultMap)
	check(err)

	var data json.RawMessage
	err = json.Unmarshal(resultMap["data"], &data)
	check(err)

	return data
}

func Fetch[ret structs.JsonStruct](method, endpoint, cookies string, headers map[string]string, reader io.Reader, complex bool, responseStruct *ret) {
	var res *http.Response
	client := &http.Client{}
	url := base_url + endpoint
	method = strings.ToUpper(method)

	if method == http.MethodGet {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		check(err)

		req.AddCookie(&http.Cookie{Name: "connect.sid", Value: cookies})
		res, err = client.Do(req)
		check(err)

		data, err := io.ReadAll(res.Body)
		check(err)

		defer res.Body.Close()
		// "result" -> "data" nesting is only for trpc endpoints
		if endpoint[1:5] == "trpc" {
			data = extractData(data)
		}

		err = json.Unmarshal(data, responseStruct)
		check(err)
		addToCache(cookies, endpoint, data)
	} else if method == http.MethodPost {
		req, err := http.NewRequest(http.MethodPost, url, reader)
		check(err)

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		req.AddCookie(&http.Cookie{Name: "connect.sid", Value: cookies})

		res, err = client.Do(req)
		check(err)

		if responseStruct != nil {
			data, err := io.ReadAll(res.Body)
			check(err)

			defer res.Body.Close()

			err = json.Unmarshal(data, responseStruct)
			check(err)
		}
	} else if method == http.MethodPut {
		req, err := http.NewRequest(http.MethodPut, url, reader)
		check(err)

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		req.AddCookie(&http.Cookie{Name: "connect.sid", Value: cookies})

		res, err = client.Do(req)
		check(err)

	}

	if res.StatusCode >= 400 {
		log.Fatalf("bad request: %d", res.StatusCode)
	}
}

func FetchTrpc[ret structs.JsonStruct](methods any, cookie string, headers map[string]string, responseStruct *ret) {
	switch m := methods.(type) {
	case []string:
		methods = strings.Join(m, ",")
	case string:
		break
	default:
		log.Fatal("invalid method type")
	}
	methods = fmt.Sprintf("/trpc/%s", methods)
	cachedData := getFromCache(cookie, methods.(string))
	if cachedData == nil {
		Fetch("get", methods.(string), cookie, headers, nil, false, responseStruct)
	} else {
		err := json.Unmarshal(cachedData, responseStruct)
		check(err)
	}
}
