package parsers

import (
	"net/url"
)

func GenericMapToQueryString(m map[string]interface{}) string {
	params := url.Values{}
	for key, value := range m {
		switch value.(type) {
		case []interface{}:
			for _, v := range value.([]interface{}) {
				sVal := InterfaceToString(v)
				params.Add(key, sVal)
			}
			continue
		}

		sVal := InterfaceToString(value)
		params.Add(key, sVal)
	}
	queryString := params.Encode()
	return queryString
}
