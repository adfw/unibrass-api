package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
	"os"
)

var DB *sql.DB

func Initialize(){
	var err error
	DB, err = sql.Open("postgres", "user=web dbname=music sslmode=disable password=PUT_PASSWORD_HERE_DO_NOT_COMMIT")
	if err != nil {
		log.Panic(err)
	}
}

func getOwner(pieceId int) (bandId int){
	rows, err := DB.Query("SELECT bandid FROM pieces WHERE pieceid = $1 LIMIT 1", pieceId)
	if err != nil {
                fmt.Fprintf(os.Stderr, "ERRRR %s", err)
		return 0
	}
	defer rows.Close()
        // There should be either 0 or 1 rows.
	for rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		if err == nil {
                        fmt.Fprintf(os.Stderr, "ERRRR %d", uid)
			bandId = uid
		} else {
			return 0
		}
	}
        fmt.Fprintf(os.Stderr, "ERRRR %d", bandId)
	return bandId
}
