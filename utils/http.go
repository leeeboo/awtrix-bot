package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func HttpPost(api string, param map[string]interface{}) ([]byte, error) {

	buf := new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(param)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Post(api, "application/json", strings.NewReader(buf.String()))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func HttpGet(api string, param map[string]interface{}) ([]byte, error) {

	queryStr, err := build(param)

	if err != nil {
		return nil, err
	}

	apiInfo, err := url.Parse(api)

	if err != nil {
		return nil, err
	}

	if apiInfo.RawQuery == "" {
		api = fmt.Sprintf("%s?%s", api, queryStr)
	} else {
		api = fmt.Sprintf("%s&%s", api, queryStr)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get(api)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
