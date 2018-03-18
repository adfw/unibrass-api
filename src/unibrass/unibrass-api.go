package main

import (
    "fmt"
    "html"
    "log"
    "net/http"
    "unibrass/database"

    "github.com/gorilla/mux"
)

func main() {

    database.Initialize()
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    log.Fatal(http.ListenAndServe(":8000", router))

}

func Index(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
