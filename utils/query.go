package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func build(raw map[string]interface{}) (string, error) {

	p := make(map[string]string, 0)

	for k, v := range raw {

		switch vv := v.(type) {
		case []interface{}:

			parseNormal(p, vv, []string{k})

			break
		case map[string]interface{}:

			parseKeyValue(p, vv, []string{k})

			break
		default:

			p[k] = fmt.Sprintf("%s", vv)

			break
		}
	}

	data := url.Values{}

	for k, v := range p {
		data.Add(k, v)
	}

	return data.Encode(), nil
}

func parseKeyValue(p map[string]string, raw map[string]interface{}, keys []string) {

	for k, v := range raw {
		switch vv := v.(type) {
		case []interface{}:

			tmpKeys := append(keys, k)

			parseNormal(p, vv, tmpKeys)

			break
		case map[string]interface{}:

			tmpKeys := append(keys, k)

			parseKeyValue(p, vv, tmpKeys)

			break
		default:

			//keys = append(keys, k)

			var tmp []string

			for m, n := range keys {
				if m > 0 {
					n = fmt.Sprintf("[%s]", n)
				}

				tmp = append(tmp, n)
			}

			kStr := strings.Join(tmp, "")

			p[fmt.Sprintf("%s[%s]", kStr, k)] = fmt.Sprintf("%s", vv)

			break
		}
	}
}

func parseNormal(p map[string]string, raw []interface{}, keys []string) {

	for k, v := range raw {
		switch vv := v.(type) {
		case []interface{}:

			tmpKeys := append(keys, fmt.Sprintf("%d", k))

			parseNormal(p, vv, tmpKeys)

			break
		case map[string]interface{}:

			tmpKeys := append(keys, fmt.Sprintf("%d", k))

			parseKeyValue(p, vv, tmpKeys)

			break
		default:

			//keys = append(keys, fmt.Sprintf("%d", k))

			var tmp []string

			for m, n := range keys {
				if m > 0 {
					n = fmt.Sprintf("[%s]", n)
				}

				tmp = append(tmp, n)
			}

			kStr := strings.Join(tmp, "")

			p[fmt.Sprintf("%s[%d]", kStr, k)] = fmt.Sprintf("%s", vv)

			break
		}
	}
}
