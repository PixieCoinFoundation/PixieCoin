package jsonloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func LoadFile(filename string, v interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("jsonloader read file error: ", filename, err.Error())
		return err
	}

	if err := json.Unmarshal(bytes, v); err != nil {
		fmt.Println("jsonloader json unmarshal error:", filename, err.Error(), "loaded string:", string(bytes))
		return err
	}

	return nil
}
