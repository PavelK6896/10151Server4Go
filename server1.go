package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"
)

type ViewData struct {
	Id      string
	Title   string
	Message string
	Users   []User
}
type User struct {
	Id    int64
	Name  string
	Email string
	State string
	Code  string
	Zip   int
}

type requestBody struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	State string `json:"state"`
	Code  string `json:"code"`
	Zip   int    `json:"zip"`
}

var countryCapitalMap map[string]string
var id int

func main() {
	countryCapitalMap := map[string]string{"France": "Paris", "Italy": "Rome", "Japan": "Tokyo", "India": "New Delhi"}
	id := 3
	fmt.Println(countryCapitalMap)
	u1 := User{1, "Bob", "b@b", "N", "432", 104}
	u2 := User{2, "Jon", "b@b", "M", "432", 123}
	u3 := User{3, "Fot", "b@b", "D", "432", 423}

	data := ViewData{
		Title: "Users List",
		Users: []User{u1, u2, u3},
	}

	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			fmt.Println(http.MethodGet)

			js, err := json.Marshal(data.Users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(js)

		} else if r.Method == http.MethodPost {
			fmt.Println(http.MethodPost)

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

			var is = true
			for _, user := range data.Users {
				if user.Name == params.Name {
					is = false
				}
			}
			if is {
				i := len(data.Users)
				id++
				newU := User{int64(id), params.Name, params.Email, params.State, strconv.Itoa(i + 1 + id), id + 100}
				data.Users = append(data.Users, newU)
			}

			if responseJson(w, err, data) {
				return
			}

		} else if r.Method == http.MethodDelete {
			fmt.Println(http.MethodDelete)

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

			for i, user := range data.Users {
				if user.Id == params.Id {
					data.Users = remove(data.Users, i)
				}
			}

			if responseJson(w, err, data) {
				return
			}
		} else if r.Method == http.MethodPut {
			fmt.Println(http.MethodPut)

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

			for i, user := range data.Users {
				if user.Id == params.Id {
					data.Users = remove(data.Users, i)
					newU := User{user.Id, params.Name, params.Email, params.State, params.Code, params.Zip}
					data.Users = append(data.Users, newU)
				}
			}

			if responseJson(w, err, data) {
				return
			}

		} else if r.Method == http.MethodPatch {
			fmt.Println(http.MethodPatch)

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

			for i, user := range data.Users {
				if user.Id == params.Id {
					data.Users[i].Zip = params.Zip
				}
			}

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

func remove(s []User, i int) []User {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func responseJson(w http.ResponseWriter, err error, data ViewData) bool {
	sort.SliceStable(data.Users, func(i, j int) bool {
		//parseInt1, _ := strconv.ParseInt(data.Users[i].Id, 10, 32)
		//parseInt2, _ := strconv.ParseInt(data.Users[j].Id, 10, 32)

		return data.Users[i].Id < data.Users[j].Id
	})

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
