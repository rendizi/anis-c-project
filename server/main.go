package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MadAppGang/httplog"
	"net/http"
	"os"
)

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	// Add other fields from the response if needed
}

var (
	loginHandler  http.Handler = http.HandlerFunc(login)
	gradesHandler http.Handler = http.HandlerFunc(Grades)
	addHandler    http.Handler = http.HandlerFunc(Additional)
	virusHandler  http.Handler = http.HandlerFunc(VirusHandler)
)

func login(w http.ResponseWriter, r *http.Request) {
	log := r.URL.Query().Get("login")
	if log == "" {
		fmt.Fprintln(w, "Login is not provided")
		return
	}
	password := r.URL.Query().Get("password")
	if password == "" {
		fmt.Fprintln(w, "Password is not provided")
		return
	}

	payload := map[string]interface{}{
		"action":      "v1/Users/Authenticate",
		"operationId": "03b5c3d1-bb48-4e97-a378-62ecf66d6c11",
		"username":    log,
		"password":    password,
		"deviceInfo":  "SM-G977N",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the request
	req, err := http.NewRequest("POST", "https://identity.micros.nis.edu.kz/v1/Users/Authenticate", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "Culture=ru-RU; lang=ru-RU")
	req.Header.Set("Host", "identity.micros.nis.edu.kz")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	fmt.Println(resp.Status)
	if resp.Status != "200 OK" {
		http.Error(w, "something", http.StatusBadRequest)
		return
	}

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Print the accessToken
	fmt.Fprintf(w, authResponse.AccessToken)
}

func main() {
	mux := http.NewServeMux()

	loggerWithFormatter := httplog.LoggerWithFormatter(httplog.DefaultLogFormatterWithRequestHeader)
	mux.Handle("/login", loggerWithFormatter(loginHandler))
	mux.Handle("/grades", gradesHandler)
	mux.Handle("/add", addHandler)
	mux.Handle("/virus", loggerWithFormatter(virusHandler))

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	err := http.ListenAndServe(":8080", corsHandler(mux))
	if err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("server closed")
		} else {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}
}
