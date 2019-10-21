package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Time struct {
	Unixtime int64 `json: "unixtime"`
}

type ChannelData struct {
	requestIP       string
	lastFetchedTime int64
	requestTime     int64
}

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Println(e)
	}
}

var startTime int64
var timeStore int64
var first bool
var counter uint32
var message chan ChannelData

func main() {
	first = true
	message = make(chan ChannelData)
	go getTime()
	go logger(message)
	http.HandleFunc("/", mainFunc)
	http.ListenAndServe(":12345", nil)
}

func getTime() {
	for {
		resp, err := http.Get("http://worldtimeapi.org/api/ip")
		// response sometimes returns "Forbidden"
		check(err)
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		var temp Time
		err = decoder.Decode(&temp)
		if err == nil {
			timeStore = temp.Unixtime
			if first == true {
				startTime = temp.Unixtime
				first = false
			}
		} else {
			check(err)
		}
		counter++
		time.Sleep(time.Second * time.Duration(math.Exp(1)))
	}
}

func mainFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Last Fetched Time: " + strconv.Itoa(int(timeStore)) + "\n"))
		w.Write([]byte("Start Time: " + strconv.Itoa(int(startTime)) + "\n"))
		w.Write([]byte("Requests: " + strconv.Itoa(int(counter))))
		message <- ChannelData{requestIP: r.RemoteAddr, lastFetchedTime: int64(timeStore), requestTime: time.Now().Unix()}
	default:
		fmt.Println("not supported")
	}
}

func logger(message chan ChannelData) {
	for {
		f, err := os.OpenFile("logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err)
		defer f.Close()
		temp := <-message
		fmt.Fprintf(f, "%s-%d-%d\n", temp.requestIP, temp.lastFetchedTime, temp.requestTime)
	}
}
