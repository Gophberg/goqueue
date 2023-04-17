package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Store struct {
	qMap  map[string][]string
	data  chan string
	topic []string
}

func (q *Store) proceed(w http.ResponseWriter, r *http.Request) {

	queueName, _ := strings.CutPrefix(r.URL.Path, "/")
	queryVal := strings.Split(r.URL.RawQuery, "=")

	switch method := r.Method; method {
	case "PUT":
		// return 400 if v is omit
		queryVal := strings.Split(r.URL.RawQuery, "=")
		if queryVal[0] != "v" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		qVal, exist := q.qMap[queueName]
		switch exist {
		case true:
			// Append val to existing key-val
			for _, e := range q.topic {
				if e == queryVal[1] {
					q.data <- queryVal[1]
				}
			}
			q.qMap[queueName] = append(qVal, queryVal[1])
		case false:
			// Creating key-val if they are not exist
			for _, e := range q.topic {
				if e == queryVal[1] {
					q.data <- queryVal[1]
				}
			}
			q.qMap[queueName] = []string{queryVal[1]}
		}
		//fmt.Println(q) // Don't forget to erase me

	case "GET":
		switch queryVal[0] {
		case "": // Return immediately
			v := pop(q, queueName)
			if v == "" {
				http.Error(w, "", http.StatusNotFound)
			}
			_, err := io.WriteString(w, v)
			if err != nil {
				fmt.Println(err)
				return
			}

		case "timeout": // Return with timeout
			v := pop(q, queueName)
			if v == "" {
				http.Error(w, "", http.StatusNotFound)
			}
			_, err := io.WriteString(w, v)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Subscribe to topic
			q.topic = append(q.topic, queryVal[0])

			t, err := strconv.Atoi(queryVal[1])
			dur := time.Second * time.Duration(t)
			if err != nil {
				panic(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), dur)

			defer cancel()

			go func(c context.Context) {

				select {
				case val := <-q.data:
					_, err := io.WriteString(w, val)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						fmt.Println(err)
						return
					}
				}
			}(ctx)

		}
	}
}

func pop(q *Store, queueName string) string {
	qVal, exist := q.qMap[queueName]
	if exist { // Response 404 if qVal is empty
		if len(qVal) == 0 {
			return ""
		}
		// Response existed requested val
		q.qMap[queueName] = qVal[1:] // Delete from qMap
		return qVal[0]
	}
	return ""
}

func NewStore() *Store {
	return &Store{
		qMap:  make(map[string][]string),
		data:  make(chan string),
		topic: make([]string, 0),
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
