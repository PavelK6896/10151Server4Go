package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type ViewData struct {
	Title   string
	Message string
	Users   []string
}

func main() {
	data := ViewData{
		Title: "Users List",
		Users: []string{"Jon", "Bob"},
	}

	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			fmt.Println("GET")

			js, err := json.Marshal(data.Users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(js)

		} else if r.Method == "POST" {
			fmt.Println("POST")

			type requestBody struct {
				Name string `json:"name"`
			}

			dat, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			params := requestBody{}
			err2 := json.Unmarshal(dat, &params)
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(params)

			data.Users = append(data.Users, params.Name)

			if responseJson(w, err, data) {
				return
			}

		} else if r.Method == "DELETE" {
			fmt.Println("DELETE")

			type requestBody struct {
				Name string `json:"name"`
			}

			dat, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			params := requestBody{}
			err2 := json.Unmarshal(dat, &params)
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(params)

			data.Users = remove(data.Users, len(data.Users)-1)

			if responseJson(w, err, data) {
				return
			}
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("templates/index1.html")
		tmpl.Execute(w, data)
	})

	fmt.Println("Server start!")
	http.ListenAndServe(":8080", nil)
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func responseJson(w http.ResponseWriter, err error, data ViewData) bool {
	js, err := json.Marshal(data.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
	return false
}
