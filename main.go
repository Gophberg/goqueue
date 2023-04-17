package main

import (
	"flag"
	"fmt"
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
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
	//	qMap.Val = qMap.Val[1:] // Delete from qMap
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
