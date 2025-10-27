package main

import (
	"encoding/json"
	"fmt"
)

type Mahasiswa struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	NIM    string `json:"nim"`
	Gender string `json:"gender"`
}

func main() {
	var object = []Mahasiswa{
		{1, "Abdul", "121212", "Male"},
		{2, "Abdi", "131313", "Male"},
		{3, "Lulu", "141414", "Female"},
	}

	jsonData, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonString := string(jsonData)
	fmt.Println(jsonString)
}
