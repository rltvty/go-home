package testhelpers

import (
	"encoding/json"
)

//IsJSON tests if a string is valid JSON
func IsJSON(str string) bool {
	var jsonStr map[string]interface{}
	err := json.Unmarshal([]byte(str), &jsonStr)
	return err == nil
}

//UnmarshalJSON Returns generic map of json data, ignoring all errors
func UnmarshalJSON(str string) map[string]interface{} {
	var jsonStr map[string]interface{}
	json.Unmarshal([]byte(str), &jsonStr)
	return jsonStr
}
