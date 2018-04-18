package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"unibrass/database"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func PieceIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Piece View")

}

func PieceView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pieceId := vars["pieceId"]
	query := `
		SELECT bands.bandid, bands.name, title, composer, arranger, notes, pieceid 
		FROM pieces
		INNER JOIN bands USING (bandid)
		WHERE pieceid = $1
		LIMIT 1`
	rows, err := database.DB.Query(query, pieceId)
	if err != nil {
		fmt.Fprintf(w, "Fatal error: %s", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := Piece{}
		err = rows.Scan(&p.Band.BandId, &p.Band.BandName, &p.Piece.Title, &p.Piece.Composer, &p.Piece.Arranger, &p.Piece.Notes, &p.PieceId)
		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(p)
		} else {
			fmt.Fprintf(w, "Fatal error: %s", err)
		}
	}

}

func PieceSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pieceName := vars["pieceName"]
	query := `
		SELECT bands.bandid, bands.name, title, composer, arranger, publisher, year, notes, pieceid 
		FROM pieces
		INNER JOIN bands USING (bandid)
		WHERE title ILIKE '%' || $1 || '%'`
	rows, err := database.DB.Query(query, pieceName)
	if err != nil {
		fmt.Fprintf(w, "Fatal error: %s", err)
		return
	}
	var pieceList PieceList
	defer rows.Close()
	for rows.Next() {
		p := Piece{}
		err = rows.Scan(&p.Band.BandId, &p.Band.BandName, &p.Piece.Title, &p.Piece.Composer, &p.Piece.Arranger, &p.Piece.Publisher, &p.Piece.Year, &p.Piece.Notes, &p.PieceId)
		if err == nil {
			pieceList = append(pieceList, p)
		} else {
			fmt.Fprintf(w, "Fatal error: %s", err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pieceList)
}

func PieceAdd(w http.ResponseWriter, r *http.Request) {
	var piece PieceUpdate
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &piece); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	query = `INSERT INTO pieces 
		(bandid, title, composer, arranger, publisher, year, notes, dateadded)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	_, dberr := database.DB.Query(
		query,
		piece.BandId,
		piece.Piece.Title,
		piece.Piece.Composer,
		piece.Piece.Arranger,
		piece.Piece.Publisher,
		piece.Piece.Year,
		piece.Piece.Notes)
	if dberr == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprintf(w, "Fatal error: %s", dberr)
		return
	}
}

func OutView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	outId := vars["outId"]
	query = `SELECT bands.bandid,
		bands.name,
		title,
		composer,
		arranger,
		notes,
		pieceid,
		outid,
		timein,
		timeout
		FROM pieces_out
		INNER JOIN pieces USING (pieceid)
		INNER JOIN bands USING (bandid)
		WHERE outid = $1
		LIMIT 1`
	rows, err := database.DB.Query(query, outId)
	if err != nil {
		fmt.Fprintf(w, "Fatal error: %s", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := Out{}
		err = rows.Scan(
			&p.Piece.Band.BandId,
			&p.Piece.Band.BandName,
			&p.Piece.Piece.Title,
			&p.Piece.Piece.Composer,
			&p.Piece.Piece.Arranger,
			&p.Piece.Piece.Notes,
			&p.Piece.PieceId,
			&p.OutId,
			&p.TimeIn,
			&p.TimeOut)
		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(p)
		} else {
			fmt.Fprintf(w, "Fatal error: %s", err)
		}
	}

}

func OutHandIn(w http.ResponseWriter, r *http.Request) {
	var outId OutInId
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &outId); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	_, dberr := database.DB.Query("UPDATE pieces_out SET timein = NOW() WHERE outid = $1", outId.OutId)
	if dberr == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprintf(w, "Fatal error: %s", dberr)
		return
	}
}

func OutHandOut(w http.ResponseWriter, r *http.Request) {
	var outId OutInId
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &outId); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	_, dberr := database.DB.Query("INSERT INTO pieces_out (pieceid) VALUES ($1)", outId.OutId)
	if dberr == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprintf(w, "Fatal error: %s", dberr)
		return
	}
}

func LoanSubmit(w http.ResponseWriter, r *http.Request) {
	var request LoanRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	_, dberr := database.DB.Query("INSERT INTO loans (pieceid, requested_by, \"from\", status) VALUES ($1, $2, NOW(), 'r')", request.PieceId, request.Requestor)
	if dberr == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprintf(w, "Fatal error: %s", dberr)
		return
	}
}

func LoanApprove(w http.ResponseWriter, r *http.Request) {
	var request LoanId
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	_, dberr := database.DB.Query("UPDATE LOANS SET status = 'a' WHERE loanid = $1", request.LoanId)
	if dberr == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprintf(w, "Fatal error: %s", dberr)
		return
	}
}
