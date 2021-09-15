package main

import (
	"fmt"
	"net/http"
	"time"
	"os"
	"bufio"
	"strings"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Hello world")
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm {
			fmt.Println(key, "=>", value)
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t := time.Now()
		result := fmt.Sprintf("%dh%d", t.Hour(), t.Minute())
		fmt.Fprintf(w, result)
	}
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case http.MethodPost:
			if err := req.ParseForm(); err != nil {
				fmt.Println("Something went bad")
				fmt.Fprintln(w, "Something went bad")
				return
			}
			
			saveData(req.FormValue("author"), req.FormValue("entry"))

			fmt.Fprintln(w, req.FormValue("author"), ": ", req.FormValue("entry"))
	}
}

func entriesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case http.MethodGet:
			saveData, err := os.ReadFile("./save.data")
			if err != nil {
				fmt.Println("Error")
				return
			}

			returnValue := ""

			values := strings.Split(strings.Trim(string(saveData), "\n") ,"\n")

			for _, value := range values {
				returnValue = returnValue + strings.Split(value, ":")[1] + "\n"
			}

			fmt.Fprintln(w, returnValue)
	}
}

func saveData(author string, entry string) {
	saveFile, err := os.OpenFile("./save.data", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	defer saveFile.Close()
	
	w := bufio.NewWriter(saveFile)

	if err == nil {
		fmt.Fprintf(w, "%s:%s\n", author, entry)
	}
	
	w.Flush()
}


func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/entries", entriesHandler)
	http.ListenAndServe(":9000", nil)
}