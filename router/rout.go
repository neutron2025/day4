package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ReadBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	page := vars["page"]
	fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", ReadBook).Methods("POST")

	// r.HandleFunc("/books/{title}", CreateBook).Methods("POST")
	// r.HandleFunc("/books/{title}", ReadBook).Methods("GET")
	// r.HandleFunc("/books/{title}", UpdateBook).Methods("PUT")
	// r.HandleFunc("/books/{title}", DeleteBook).Methods("DELETE")

	http.ListenAndServe(":80", r)
}

//go mod init <module-name>
// go get -u github.com/gorilla/mux
