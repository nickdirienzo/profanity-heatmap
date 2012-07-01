package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	//"net/http"
	"os"
)

type Commit struct {
	Payload Payload `json:"payload"`
	Event   string  `json:"type"`
}

type Payload struct {
	Shas []map[string]interface{} `json:"shas"`
}

func main() {
	/*resp, err := http.Get("http://data.githubarchive.org/2012-03-11-12.json.gz")
	if err != nil {
		fmt.Println("Couldn't retrieve data from GitHub Archive")
		return
	}
	defer resp.Body.Close()
	reader, err := gzip.NewReader(resp.Body)*/
	file, err := os.Open("2012-03-11-15.json.gz")
	if err != nil {
		fmt.Println("Couldn't open file")
		return
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println("Couldn't uncompress")
		return
	}
	decoder := json.NewDecoder(reader)
	for {
		var commit Commit
		if err := decoder.Decode(&commit); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Couldn't decode JSON:", err)
		}
		if commit.Event == "PushEvent" {
			fmt.Println(commit.Payload.Shas)
		}
	}
}
