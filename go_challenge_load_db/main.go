package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strings"
)

func main() {
	location := os.Args[1]
	fmt.Println(location)
	loadDatabase(location)

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
