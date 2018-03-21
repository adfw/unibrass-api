/*****************************************************************************
 * UniBrass Library API                                                      *
 *                                                                           *
 * Copyright (C) Anthony Williams, 2018                                      *
 *****************************************************************************/

package main

import (
	"log"
	"net/http"

	"unibrass/database"
)

func main() {

	database.Initialize()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))

}
