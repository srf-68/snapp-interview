package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"log"
	"time"
	"regexp"
	"strings"
	"strconv"
  
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "snapp"
)  

func indexQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	rgx := regexp.MustCompile("[^\\w(?.)]+|([,\\/#!$%\\^&\\*;:{}=\\-_\\?~()])")
	splitted := rgx.Split(query, -1)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " + "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	sqlStatement := `INSERT INTO queries (term, millitime) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, query, time.Now().UnixNano())

	for _, element := range splitted {
		if len(element) > 0 {
			temp := element
			if strings.HasSuffix(element, ".") {
				temp = element[:len(element)-1]
			}
			if strings.HasPrefix(element, ".") {
				temp = element[1:len(element)]
			}
			_, err = db.Exec(sqlStatement, temp, time.Now().UnixNano())
		}
	}

	if err != nil {
  		panic(err)
	}
	db.Close()
}

func returnQueries(w http.ResponseWriter, r *http.Request) {
	hour := r.URL.Query().Get("hour")
	count := r.URL.Query().Get("count")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " + "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	i, err := strconv.ParseInt(hour, 10, 64)
	nanotime := time.Now().UnixNano() - i * 60 * 60 * 1000000000
	rows, err := db.Query("SELECT term, count(term) as cnt FROM queries WHERE millitime > $1 GROUP BY term LIMIT $2", nanotime, count)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	var csv string
	for rows.Next() {
		var term string
		var cnt int
		err = rows.Scan(&term, &cnt)
		if err != nil {
			// handle this error
			panic(err)
		}
		csv += fmt.Sprintf("%v,%v\n", term, cnt)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	db.Close()
	w.Header().Set("Content-Type", "text/csv")
    w.WriteHeader(200)
    w.Write([]byte(csv))
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/index-search-query", indexQuery)
	router.HandleFunc("/return-queries", returnQueries)
	log.Fatal(http.ListenAndServe(":8585", router))
}