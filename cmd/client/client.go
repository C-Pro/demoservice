package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"demoservice/pkg/service"
)

const (
	host   = "localhost"
	port   = 8080
	schema = "http"
	method = "time"
)

func buildURL() string {
	return fmt.Sprintf("%s://%s:%d/%s", schema, host, port, method)
}

func main() {
	client := http.DefaultClient
	resp, err := client.Get(buildURL())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ts := service.TimeStruct{}

	if err := json.Unmarshal(b, &ts); err != nil {
		panic(err)
	}

	fmt.Println(ts.Date)
	fmt.Println(ts.Time)
}
