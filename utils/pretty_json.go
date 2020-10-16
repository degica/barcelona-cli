package utils

import (
	"encoding/json"
)

func PrettyJSON(b []byte) string {
	var v interface{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return ""
	}
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(pretty)
}
