package main

import (
	"database/sql"

	"github.com/lib/pq"
)

type Band struct {
	BandId		int	`json:"bandid"`
	BandName	string	`json:"bandname"`
}

type SQLPiece struct {
	Title		sql.NullString	`json:"title"`
	Composer	sql.NullString	`json:"composer"`
	Arranger	sql.NullString	`json:"arranger"`
	Publisher	sql.NullString	`json:"publisher"`
	Year		pq.NullTime	`json:"year"`
	Notes		sql.NullString	`json:"notes"`
}

type Piece struct {
	PieceId		int		`json:"pieceid"`
	Band		Band		`json:"band"`
	Piece		SQLPiece	`json:"piece"`
}

type PieceUpdate struct {
	Piece	SQLPiece	`json:"piece"`
	BandId	int		`json:"bandid"`
}

type PieceList []Piece

type Loan struct {
	LoanId		int		`json:"loanid"`
	PieceId		int		`json:"pieceid"`
	Lender		int		`json:"lender"`
	Requestor	int		`json:"requestor"`
	Status		string		`json:"status"`
	DateFrom	pq.NullTime	`json:"datefrom"`
	DateDue		pq.NullTime	`json:"dateuntil"`
	DateSent	pq.NullTime	`json:"datesent"`
	DateReturned	pq.NullTime	`json:"datereturned"`
}

type LoanList []Loan

type Out struct {
	OutId   int         `json:"outid"`
	TimeOut pq.NullTime `json:"timeout"`
	TimeIn  pq.NullTime `json:"timein"`
	Piece   Piece       `json:"piece"`
}

type OutInId struct {
	OutId int `json:"outid"`
}
