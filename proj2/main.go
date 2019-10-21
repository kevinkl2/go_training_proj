package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Time struct {
	DateTime string `json: "datetime"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var startTime string
var timeStore string
var first bool
var counter uint32

func main() {
	first = true
	go getTime()
	http.HandleFunc("/", mainFunc)
	http.ListenAndServe(":12345", nil)
}

func getTime() {
	for {
		resp, err := http.Get("http://worldtimeapi.org/api/ip")
		check(err)
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		var temp Time
		err = decoder.Decode(&temp)
		timeStore = temp.DateTime
		check(err)
		if first == true {
			startTime = temp.DateTime
			first = false
		}
		counter++
		time.Sleep(time.Second * time.Duration(math.Exp(1)))
	}
}

func mainFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Last Fetched Time: " + timeStore + "\n"))
		w.Write([]byte("Start Time: " + startTime + "\n"))
		w.Write([]byte("Requests: " + strconv.Itoa(int(counter))))
	default:
		fmt.Println("not supported")
	}
}
