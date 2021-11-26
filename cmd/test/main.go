package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Minimum int    `json:"minimum"`
	Unit    string `json:"unit"`
}

func main() {
	var str = `{
          "minimum": 12,
          "unit": "milliseconds"
        }`
	var i interface{}
	err := json.Unmarshal([]byte(str), &i)
	if err != nil {
		fmt.Println("1111111111" + err.Error())
	}
	var m Message
	data, err := json.Marshal(i)
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println(err.Error())
	}

}
