package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Values struct {
	Source       int `json:"source"`
	Environments []struct {
		Name      string         `json:"name"`
		Variables map[string]any `json:"variables"`
	} `json:"environments"`
}

func main() {
	url := os.Args[1]
	fmt.Println(url)
	UrlRequest(url)

}

func UrlRequest(url string) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//fmt.Println(body)
	response := UnmJson(body)
	PrintData(response)
}

func UnmJson(body []byte) Values {
	var datafile Values
	err := json.Unmarshal(body, &datafile)
	if err != nil {
		log.Fatal(err)

	}
	return datafile

}

func PrintData(datafile Values) {
	for _, value := range datafile.Environments {
		fmt.Println(value.Name)
		for z, a := range value.Variables {
			for _, j := range a.(map[string]interface{}) {
				d2 := z + ": " + j.(string)
				fmt.Println(d2)

			}

		}

	}

}
