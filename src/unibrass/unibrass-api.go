/*****************************************************************************
 * UniBrass Library API                                                      *
 *                                                                           *
 * Copyright (C) Anthony Williams, 2018                                      *
 *****************************************************************************/

package main

import (

	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
    "database/sql"

	"unibrass/database"

    "github.com/gorilla/mux"
)

type Piece struct {
	Band     string			`json:"owner"`
	PieceId  int			`json:"pieceid"`
	Title    sql.NullString	`json:"title"`
	Composer sql.NullString	`json:"composer"`
	Arranger sql.NullString `json:"arranger"`
	Notes    sql.NullString `json:"notes"`
}

func main() {

	database.Initialize()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/piece", PieceIndex)
	router.HandleFunc("/piece/{pieceId}", PieceView)
	log.Fatal(http.ListenAndServe(":8000", router))

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func PieceIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Piece View")

}

func PieceView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pieceId := vars["pieceId"]
	rows, err := database.DB.Query("SELECT bands.name, title, composer, arranger, notes, pieceid FROM pieces INNER JOIN bands USING (bandid) WHERE pieceid = $1 LIMIT 1", pieceId)
	if err != nil {
		fmt.Fprintf(w, "Fatal error: %s", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := Piece{}
		err = rows.Scan(&p.Band, &p.Title, &p.Composer, &p.Arranger, &p.Notes, &p.PieceId)
		if err == nil {
			json.NewEncoder(w).Encode(p)
		} else {
			fmt.Fprintf(w, "Fatal error: %s", err)
		}
	}

}
