package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
)

type Queue struct {
	Name string   `json:"name"`
	Val  []string `json:"val"`
}

func (q Queue) proceed(w http.ResponseWriter, r *http.Request) {

	path, _ := strings.CutPrefix(r.URL.Path, "/")

	switch rawQeury := strings.Split(r.URL.RawQuery, "="); rawQeury[0] {
	case "v":
		q.Name = path
		q.Val = append(q.Val, rawQeury[1])

		fmt.Println(q)
	case "timeout":
		fmt.Println(rawQeury[1])
	default:
		http.Error(w, "", http.StatusBadRequest)
	}

	//	q.Val = q.Val[1:] // Delete from queue
}

func main() {
	var port int
	flag.IntVar(&port, "port", 3000, "set the port for http server without colon, like `3000`")
	flag.Parse()

	q := Queue{}

	http.HandleFunc("/", q.proceed)

	sPort := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting at %s\n", sPort)
	err := http.ListenAndServe(sPort, nil)
	if err != nil {
		panic(err)
	}
}
