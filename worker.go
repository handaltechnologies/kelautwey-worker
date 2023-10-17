package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

var url string

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func checkStatus(url string, ch chan<- string) {
	response, err := http.Get(url)
	if err != nil {
		ch <- url + " " + err.Error()
		return
	}
	defer response.Body.Close()
	ch <- url + " " + response.Status
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	rows, err := db.Query("SELECT url FROM websites")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	ch := make(chan string)
	urls := 0
	for rows.Next() {
		err := rows.Scan(&url)
		if err != nil {
			panic(err)
		}
		//fmt.Println(url)
		go checkStatus(url, ch)
		urls++
	}

	for i := 0; i <= urls; i++ {
		fmt.Println(<-ch)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

}
