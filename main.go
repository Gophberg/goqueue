package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Store struct {
	qMap map[string][]string
}

func (q *Store) proceed(w http.ResponseWriter, r *http.Request) {

	qName, _ := strings.CutPrefix(r.URL.Path, "/")

	switch queryVal := strings.Split(r.URL.RawQuery, "="); queryVal[0] {
	case "v":

		qVal, exist := q.qMap[qName]
		switch exist {
		case true:
			// Append val to existing key-val
			q.qMap[qName] = append(qVal, queryVal[1])
		case false:
			// Creating key-val if they are not exist
			q.qMap[qName] = []string{queryVal[1]}
		}

		fmt.Println(q)

	case "timeout":
		fmt.Println(queryVal[1])

	case "": // Return val to requested key
		qVal, exist := q.qMap[qName]
		switch exist {
		case true:
			// Response 404 if qVal is empty
			if len(qVal) == 0 {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			// Response existed requested val
			_, err := io.WriteString(w, qVal[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			q.qMap[qName] = qVal[1:] // Delete from qMap
		case false: // Response 404 if key-val was not created
			http.Error(w, "", http.StatusNotFound)
		}

	default: // Return 400 if v body is empty
		http.Error(w, "", http.StatusBadRequest)
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
