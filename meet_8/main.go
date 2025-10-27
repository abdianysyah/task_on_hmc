package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Id		int 	`json:id`
	Name 	string 	`json:name`
	NIM 	string	`json:nim`
	Gender	string	`json:gender`  
}


func main()  {
	// Decode  Json
	var jsonString = `{"id" : 1, "name" : "Abdul", "nim" : "121212", "gender" : "male"}`
	var jsonData = []byte(jsonString)

	var data User
	var err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Id	:", data.Id)
	fmt.Println("Name	:", data.Name)
	fmt.Println("NIM	:", data.NIM)
	fmt.Println("Gender	:", data.Gender)
}

