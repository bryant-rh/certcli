package util

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/logrusorgru/aurora"
)

// ContainsInSlice 判断字符串是否在 slice 中
func ContainsInSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func JsonToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

func RedColor(s interface{}) string {
	return fmt.Sprintf("%s", aurora.Red(s))
}

func YellowColor(s interface{}) string {
	return fmt.Sprintf("%s", aurora.Yellow(s))
}

func GreenColor(s interface{}) string {
	return fmt.Sprintf("%s", aurora.Green(s))
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

//MapKeyInIntSlice
func MapKeyInIntSlice(haystack []string, needle string) bool {
	set := make(map[string]struct{})

	for _, e := range haystack {
		set[e] = struct{}{}
	}

	_, ok := set[needle]
	return ok
}
