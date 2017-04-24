package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/urlinfo/1/{url:.*}", safetyCheck)
	http.Handle("/", r)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

func safetyCheck(w http.ResponseWriter, r *http.Request) {
	var url string
	w.Header().Add("Content-type", "application/json")
	vars := mux.Vars(r)
	domain := vars["url"]
	query := r.URL.RawQuery
	if query != "" {
		domain += "?" + query
	}
	db, _ := sql.Open("mysql", "go:GO@(localhost:3333)/go_challenge")
	defer db.Close()
	err := db.QueryRow("SELECT url FROM blocked_urls where url='" + domain + "';").Scan(&url)
	switch {
	case err == sql.ErrNoRows:
		w.Write([]byte("{\"safe\": true, \"message\": \"URL is safe\"}"))
	case err != nil:
		fmt.Println(err)
	default:
		w.Write([]byte("{\"safe\": false, \"message\": \"URL is not safe\"}"))
	}
	//w.Write([]byte(domain))

}
