package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type VirusResponse struct {
	Positives int    `json:"positives"`
	Total     int    `json:"total"`
	Link      string `json:"permalink"`
	Msg       string `json:"verbose_msg"`
}

var apiKey = "aead23343c61d09c8ba9623e93375126cf517bb150484653cd5ce861d0693cbd"

func VirusHandler(w http.ResponseWriter, r *http.Request) {
	sh := r.URL.Query().Get("url")
	if sh == "" {
		http.Error(w, "empty", http.StatusInternalServerError)
		return
	}
	fmt.Println("Error parsing JSON:")

	urle, err := Fromb("https://b2a.kz/" + sh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	apiURL := "http://www.virustotal.com/vtapi/v2/url/report"
	data := url.Values{}
	data.Set("apikey", apiKey)
	data.Set("resource", urle)

	reqBody := strings.NewReader(data.Encode())
	fmt.Println("Error parsing JSON:", err)

	req, err := http.NewRequest("POST", apiURL, reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Error parsing JSON:", err)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Response status is not ok", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Error parsing JSON:", err)

	var response VirusResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if response.Msg == "Resource does not exist in the dataset" {
		link, err := ScanUrl(urle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, link)
		return
	}
	fmt.Fprintln(w, response.Positives, "/", response.Total)
}

func Fromb(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Header.Get("Location"), nil
}

func ScanUrl(urle string) (string, error) {
	fmt.Println(urle)
	apiURL := "https://www.virustotal.com/vtapi/v2/url/scan"
	data := url.Values{}
	data.Set("apikey", apiKey)
	data.Set("url", urle)

	reqBody := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", apiURL, reqBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("125")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("response status is not ok")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response VirusResponse
	fmt.Println(string(body))
	err = json.Unmarshal(body, &response)

	if err != nil {
		return "", err
	}
	return response.Link, nil
}
