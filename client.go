package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Person struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

//client connection check
func main() {
	fmt.Println("Client Started")
	appPort := "8000"
	response, err := http.Get("http://localhost:" + appPort + "/get")
	if err != nil {
		fmt.Println("HTTP req failed with error", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
