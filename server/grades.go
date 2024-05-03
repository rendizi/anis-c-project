package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Grades(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is empry", http.StatusBadRequest)
		return
	}
	user, err := DecodeJWT(token)
	log.Println(user["PersonGid"].(string))
	studentId := r.URL.Query().Get("id")
	if studentId == "" {
		http.Error(w, "student id is empty", http.StatusBadRequest)
		return
	}
	payload := map[string]interface{}{
		"operationId": "9ecdebdd-697e-4d27-9cbf-f768ccf63080",
		"action":      "v1/ReportCard/GetAllReportCardsAsync",
		"studentId":   user["PersonGid"].(string),
		"token":       token,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the request
	req, err := http.NewRequest("POST", "https://reportcard.micros.nis.edu.kz/v1/ReportCard/GetAllReportCardsAsync", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "Culture=ru-RU; lang=ru-RU")
	req.Header.Set("authorization", token)
	req.Header.Set("Host", "reportcard.micros.nis.edu.kz")

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

	var data []map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	n := len(data)
	res := make([]string, n)
	for i := n - 1; i >= 0; i-- {
		cur, _ := data[i]["schoolYear"].(map[string]interface{})["name"].(map[string]interface{})["ru"].(string)
		sum := 0
		reportCards, _ := data[i]["reportCard"].([]interface{})
		for _, card := range reportCards {
			resultMark := card.(map[string]interface{})["resultMark"]
			if resultMark == nil {
				sum += 0
			} else if result, ok := resultMark.(map[string]interface{}); ok {
				if ruValue, ok := result["ru"].(string); ok {
					if ruValue == "зачет" {
						sum += 5
					} else if intValue, err := strconv.Atoi(ruValue); err == nil {
						sum += intValue
					} else {
						sum += 5
					}
				} else {
					continue
				}
			} else {
				continue
			}
		}
		s := fmt.Sprintf("%s:%.2f", cur, float64(sum)/float64(len(reportCards)))
		res[i] = s
	}
	fmt.Fprintln(w, strings.Join(res, ";"))

}

func DecodeJWT(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %v", err)
	}

	var userInfo map[string]interface{}
	if err = json.Unmarshal([]byte(claims["UserInfo"].(string)), &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return userInfo, nil
}
