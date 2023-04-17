package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Store struct {
	qMap map[string][]string
}

type Consumer struct {
	ch      chan string
	timeout time.Time
}

func (q *Store) proceed(w http.ResponseWriter, r *http.Request) {

	queueName, _ := strings.CutPrefix(r.URL.Path, "/")
	queryVal := strings.Split(r.URL.RawQuery, "=")

	switch method := r.Method; method {
	case "PUT":
		qVal, exist := q.qMap[queueName]

		switch exist {
		case true:
			// Append val to existing key-val
			println("Append val to existing key-val")
			q.qMap[queueName] = append(qVal, queryVal[1])
		case false:
			// Creating key-val if they are not exist
			println("Creating key-val if they are not exist")
			q.qMap[queueName] = []string{queryVal[1]}
		}
		// if v is omit
		queryVal := strings.Split(r.URL.RawQuery, "=")
		if queryVal[0] == "" {
			http.Error(w, "", http.StatusBadRequest)
		}
		fmt.Println(q) // Don't forget to erase me

	case "GET":
		switch queryVal[0] {
		case "": // Return immediately
			qVal, exist := q.qMap[queueName]
			switch exist {
			case true:
				// Response 404 if qVal is empty
				if len(qVal) == 0 {
					println("404 qVal is empty")
					http.Error(w, "", http.StatusNotFound)
					return
				}
				// Response existed requested val
				_, err := io.WriteString(w, qVal[0])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					fmt.Println(err)
					return
				}
				q.qMap[queueName] = qVal[1:] // Delete from qMap
			case false: // Response 404 if key-val was not created
				println("404 key-val was not created")
				http.Error(w, "", http.StatusNotFound)
			}

		case "timeout":
			fmt.Println(queryVal[1])
		}
	}
}

func NewStore() *Store {
	return &Store{
		qMap: make(map[string][]string),
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 3000, "set the port for http server without colon, like `3000`")
	flag.Parse()

	store := NewStore()

	http.HandleFunc("/", store.proceed)

	sPort := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting at %s\n", sPort)
	err := http.ListenAndServe(sPort, nil)
	if err != nil {
		panic(err)
	}
}
