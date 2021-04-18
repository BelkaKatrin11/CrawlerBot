package main

import (
	"CrawlerBot/internal/models"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const source = "https://lenta.ru/rss/news/russia"

func main() {
	response, err := getResponse()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = parseResponse(response)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getResponse() ([]byte, error) {
	resp, err := http.Get(source)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return nil, errors.New("wrong status code")
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseResponse(data []byte) error {
	var rss models.RSS

	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return err
	}

	for _, item := range rss.Channel.Item {
		fmt.Printf("%+v\n\n", item)
	}

	return nil
}