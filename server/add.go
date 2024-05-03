package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SchoolName struct {
	Kazakh  string `json:"kk"`
	Russian string `json:"ru"`
	English string `json:"en"`
}

type School struct {
	ID   string     `json:"gid"`
	Name SchoolName `json:"name"`
}

type Response struct {
	PhotoURL string        `json:"photoUrl"`
	Class    string        `json:"klass"`
	School   School        `json:"school"`
	Children []interface{} `json:"children"`
}

func Additional(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is empty", http.StatusBadRequest)
		return
	}
	payload := map[string]interface{}{
		"applicationType": "ContingentAPI",
		"action":          "Api/AdditionalUserInfo",
		"operationId":     "7ea4e116-5585-4518-b180-453815338986",
		"token":           token,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "https://contingent.micros.nis.edu.kz/Api/AdditionalUserInfo", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "Culture=ru-RU; lang=ru-RU")
	req.Header.Set("authorization", token)
	req.Header.Set("Host", "contingent.micros.nis.edu.kz")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", http.StatusText(resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Fprintln(w, response.School.ID+";"+response.Class)

}
