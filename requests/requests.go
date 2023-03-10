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

func Fetch[ret structs.JsonStruct](method, endpoint, cookies, body string, complex bool, responseStruct *ret) {
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

		bytes, err := io.ReadAll(res.Body)
		check(err)

		defer res.Body.Close()
		data := extractData(bytes)

		err = json.Unmarshal(data, responseStruct)
		check(err)
	}

	if res.StatusCode >= 400 {
		log.Fatal("bad request")
	}
}

func FetchTrpc[ret structs.JsonStruct](methods []string, cookie string, responseStruct *ret) {
	m := strings.Join(methods, ",")
	Fetch("get", fmt.Sprintf("/trpc/%s", m), cookie, "", false, responseStruct)
}