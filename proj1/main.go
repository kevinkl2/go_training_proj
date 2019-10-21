package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Person struct {
	Name       string `json: "name"`
	Age        int    `json: "age"`
	Profession string `json: "profession"`
	HairColor  string `json: "hairColor"`
}

var personMap map[string]*Person

func main() {
	personMap = make(map[string](*Person))
	http.HandleFunc("/person", personFunc)
	http.HandleFunc("/person/", searchFunc)
	http.ListenAndServe(":8080", nil)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func personFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		f, err := os.Create("test.json")
		check(err)
		defer f.Close()

		for _, value := range personMap {
			marshalled, _ := json.Marshal(value)
			fmt.Fprintln(w, string(marshalled))
			_, err = f.WriteString(string(marshalled) + "\n")
			f.Sync()
		}
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var temp Person
		err := decoder.Decode(&temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		personMap[temp.Name] = &temp
	default:
		fmt.Println("not supported")
	}
}

func searchFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		name := r.URL.Path[8:]
		// ERROR IF NAME IS NOT FOUND
		retrieved, _ := personMap[name]
		marshalled, _ := json.Marshal(*retrieved)
		fmt.Fprintf(w, string(marshalled))
	case "POST":
		fmt.Println("POST")
	default:
		fmt.Println("not supported")
	}
}
