package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
)

func main() {
	loadDatabase("../go_challenge_load_db/blacklist")
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

func loadDatabase(location string) {
	var url string
	db, _ := sql.Open("mysql", "go:GO@(localhost:3333)/go_challenge")
	defer db.Close()

	if strings.HasPrefix(location, "http") {
		resp, err := http.Get(location)
		defer resp.Body.Close()
		if err == nil {
			body := bufio.NewScanner(resp.Body)
			for body.Scan() {
				if strings.HasPrefix(body.Text(), "#") {
					continue
				}
				err := db.QueryRow("SELECT url FROM blocked_urls where url='" + body.Text() + "';").Scan(&url)
				switch {
				case err == sql.ErrNoRows:
					_, err := db.Exec("INSERT INTO blocked_urls (url) VALUES ('" + body.Text() + "' );")
					if err != nil {
						fmt.Println(err)
					}
				case err != nil:
					fmt.Println(err)
				default:
					continue
				}
			}
		}
	} else {
		file, err := os.Open(location)
		defer file.Close()
		if err == nil {
			reader := bufio.NewScanner(file)
			for reader.Scan() {
				if strings.HasPrefix(reader.Text(), "#") {
					continue
				}
				err := db.QueryRow("SELECT url FROM blocked_urls where url='" + reader.Text() + "';").Scan(&url)
				switch {
				case err == sql.ErrNoRows:
					_, err := db.Exec("INSERT INTO blocked_urls (url) VALUES ('" + reader.Text() + "' );")
					if err != nil {
						fmt.Println(err)
					}
				case err != nil:
					fmt.Println(err)
				default:
					continue
				}
			}
		} else {
			println(err.Error())
		}
	}

}