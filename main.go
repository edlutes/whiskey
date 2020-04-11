package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Whiskey struct {
	Name        string  `json:"name"`
	Price       float32 `json:"price,omitempty"`
	Location    string  `json:"location,omitempty"`
	Distillery  string  `json:"distillery,omitempty"`
	Type        string  `json:"type,omitempty"`
	Description string  `json:"description,omitempty"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	fmt.Println("Hi")
	http.HandleFunc("/hello", hello)

	connString := "dbname=whiskey_sample sslmode=disable"
	db, err := sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})

	http.HandleFunc("/get", getWhiskeyHandler)
	http.HandleFunc("/enter", enterWhiskeyHandler)

	http.ListenAndServe(":3000", nil)
}

func getWhiskeyHandler(w http.ResponseWriter, r *http.Request) {

	whiskeys, err := store.getWhiskeys()

	whiskeyListBytes, err := json.Marshal(whiskeys)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(whiskeyListBytes)
}

func enterWhiskeyHandler(w http.ResponseWriter, r *http.Request) {
	whiskey := Whiskey{}

	err := r.ParseForm()

	for key, value := range r.Form {
		fmt.Println("%s = %s", key, value)
	}

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	whiskey.Name = r.Form.Get("name")
	whiskey.Description = r.Form.Get("description")

	fmt.Println(whiskey.Name, whiskey.Description)

	err = store.enterWhiskey(&whiskey)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/assets/", http.StatusFound)
}
