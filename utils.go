package main

import (
	"encoding/json"
	"strconv"
)

// ToJSON util
func ToJSON(v interface{}, indents ...bool) string {
	var dt []byte
	if len(indents) > 0 && indents[0] {
		dt, _ = json.MarshalIndent(v, "", "  ")
	} else {
		dt, _ = json.Marshal(v)
	}

	return string(dt)
}

// ToStr util
func ToStr(i interface{}) string {
	if i == nil {
		return ""
	}
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case int8:
		return strconv.Itoa(int(i.(int8)))
	case int16:
		return strconv.Itoa(int(i.(int16)))
	case int32:
		return strconv.Itoa(int(i.(int32)))
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	}
	return ""
}

