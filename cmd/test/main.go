package main

import (
	"encoding/json"
	"fmt"
)

type String = string
type Array = []string

type ArrayOfString struct {
	String
	Array
}

func main() {
	str := `{"12123","sdf"}`
	var as ArrayOfString
	err := json.Unmarshal([]byte(str), &as)
	if err != nil {
		fmt.Println(err.Error())
	}
}
