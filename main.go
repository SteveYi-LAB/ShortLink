package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func webServer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Method:", r.Method)
	fmt.Println(r.URL.Path)
	fmt.Println("IP Address: ", getIP(r))

	p := "." + r.URL.Path

	if r.Method == "GET" {
		if p == "./" || p == "./index.html" {
			http.ServeFile(w, r, "./index.html")
		} else {
			shortLinkCode := strings.ReplaceAll(p, "./", "")
			fmt.Println("Code: " + shortLinkCode)

			db, err := sql.Open("sqlite3", "db")
			if err != nil {
				fmt.Println(err)
			}
			rows, err := db.Query("SELECT link FROM shortlink WHERE code = ?", shortLinkCode)

			for rows.Next() {
				var link string
				err = rows.Scan(&link)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Link: " + link)
				http.Redirect(w, r, link, 302)
			}
		}
	}

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if r.URL.Path == "/api/create" {
			tokenValue := "123456"

			r.ParseForm()
			token := "NULL"
			link := "NULL"

			for key, values := range r.Form {
				if key == "token" {
					token = values[0]
				}
				if key == "link" {
					link = values[0]
				}
			}

			if token == tokenValue {
				code := randomString(5)

				db, err := sql.Open("sqlite3", "db")
				if err != nil {
					fmt.Println(err)
				}

				stmt, err := db.Prepare("INSERT INTO shortlink(code, link, ipAddress) values(?,?,?)")
				if err != nil {
					fmt.Println(err)
				}
				res, err := stmt.Exec(code, link, getIP(r))
				if err != nil {
					fmt.Println(err)
				}
				id, err := res.LastInsertId()
				if err != nil {
					fmt.Println(err)
					fmt.Println(id)
				} else {
					io.WriteString(w, "Success!\n")
					io.WriteString(w, code)
				}
			} else {
				io.WriteString(w, "Unauthorized!\n")
			}
		}
	}
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func randomString(length int) string {
	rand.Seed(time.Now().Unix())

	var output strings.Builder

	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return (output.String())
}

func main() {
	http.HandleFunc("/", webServer)

	fmt.Print("\n")
	fmt.Print("-------------------\n")
	fmt.Print("\n")
	fmt.Print("Port listing at 32156\n")
	fmt.Print("Author: SteveYi\n")
	fmt.Print("Demo: https://yiy.tw/\n")
	fmt.Print("\n")
	fmt.Print("-------------------\n")
	fmt.Print("\n")

	err := http.ListenAndServe(":32156", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
