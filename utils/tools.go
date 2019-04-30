package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"reflect"
	"strings"
)

func ParseMapToStruct(raw interface{}, dst interface{}) error {

	jsonData, err := json.Marshal(raw)

	if err != nil {
		return errors.New("[ParseMapToStruct] " + err.Error())
	}

	d := json.NewDecoder(strings.NewReader(string(jsonData)))
	d.UseNumber()
	err = d.Decode(dst)
	if err != nil {
		return errors.New("[ParseMapToStruct] " + err.Error())
	}
	return nil
}

func Md5(raw string) string {

	m := md5.New()
	m.Write([]byte(raw))
	return hex.EncodeToString(m.Sum(nil))
}

func GetIP(req *http.Request) (string, error) {

	if req == nil {
		return "", errors.New("Request is nil.")
	}

	forward := req.Header.Get("X-Forwarded-For")

	if forward != "" {
		return forward, nil
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", errors.New("IP parse error.")
	}

	return ip, nil
}

func InArray(needle interface{}, haystack interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}
