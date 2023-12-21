package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	var data *map[string]any
	fmt.Println("data: ")
	fmt.Println(data)

	data = &map[string]any{}
	(*data)["test1"] = "test1"
	(*data)["test2"] = 2
	(*data)["test3"] = nil
	// (*data)["test4"] =
	(*data)["test5"] = data

	fmt.Println("data: ")
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(buf)
	}

	testValue1, ok := (*data)["test1"].(string)
	fmt.Println("test1, ok: ")
	fmt.Println(testValue1 + ", " + strconv.FormatBool(ok))

	testValue2, ok := (*data)["test2"].(string)
	fmt.Println("test2, ok: ")
	fmt.Println(testValue2 + ", " + strconv.FormatBool(ok))

	testValue3, ok := (*data)["test3"].(string)
	fmt.Println("test3, ok: ")
	fmt.Println(testValue3 + ", " + strconv.FormatBool(ok))

	testValue4, ok := (*data)["test4"].(string)
	fmt.Println("test4, ok: ")
	fmt.Println(testValue4 + ", " + strconv.FormatBool(ok))

	testValue5, ok := (*data)["test5"].(string)
	fmt.Println("test5, ok: ")
	fmt.Println(testValue5 + ", " + strconv.FormatBool(ok))
}
