package main

import (
	"database/sql"
	"github.com/lib/pq"
)

type Piece struct {
	Band     string         `json:"owner"`
	PieceId  int            `json:"pieceid"`
	Title    sql.NullString `json:"title"`
	Composer sql.NullString `json:"composer"`
	Arranger sql.NullString `json:"arranger"`
	Notes    sql.NullString `json:"notes"`
}

type Out struct {
	OutId   int         `json:"outid"`
	TimeOut pq.NullTime `json:"timeout"`
	TimeIn  pq.NullTime `json:"timein"`
	Piece   Piece       `json:"piece"`
}

type OutInId struct {
	OutId int `json:"outid"`
}
