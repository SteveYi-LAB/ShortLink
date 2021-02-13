package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func webServer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\nMethod:", r.Method)
	fmt.Println(r.URL.Path)
	fmt.Println("IP Address: ", getIP(r))
	fmt.Println("Current Time: ", time.Now())

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

			var linkcheck int

			for rows.Next() {
				var link string
				err = rows.Scan(&link)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Link: " + link)
				if link == "" {
					linkcheck = 0
				} else {
					http.Redirect(w, r, link, 302)
					linkcheck = 1
				}
			}

			if linkcheck == 0 {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, "Link not found!\n")
				fmt.Println("Not Found!")
			}
		}
	}

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if r.URL.Path == "/api/create" {
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}

			tokenValue := os.Getenv("token")

			r.ParseForm()
			token := "NULL"
			link := "NULL"
			admin := "false"
			googleRecaptcha := "NULL"
			custom := "false"
			customcode := ""

			for key, values := range r.Form {
				if key == "token" {
					token = values[0]
				}
				if key == "link" {
					link = values[0]
				}
				if key == "admin" {
					admin = values[0]
				}
				if key == "g-recaptcha-response" {
					googleRecaptcha = values[0]
				}
				if key == "custom" {
					custom = values[0]
				}
				if key == "customcode" {
					customcode = values[0]
				}
			}

			if strings.Contains(link, "http://") && strings.Contains(link, "https://") {
				if admin == "true" {
					if token == tokenValue {

						var code string
						code = randomString(5)

						if custom == "true" {
							code = customcode
						}

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
				} else {
					if verifyRecaptcha(googleRecaptcha) == "1" {
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
							io.WriteString(w, "Code: "+code+"\n")
							io.WriteString(w, "https://yiy.tw/"+code+"\n")
							fmt.Println("Link Create!")
							fmt.Println("Code: " + code)
							fmt.Println("Link: " + link)
						}
					} else {
						io.WriteString(w, "Google Recaptcha Failure!\n")
					}
				}
			} else {
				io.WriteString(w, "Please type a Vaild URL.")
			}
		} else {
			io.WriteString(w, "HTTP "+r.Method+" Not Support!")
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

func verifyRecaptcha(recaptcha string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	verifyLink := "https://www.recaptcha.net/recaptcha/api/siteverify"
	secretKey := os.Getenv("Google_Recaptcha_SecretKey")
	verifySuccess := "0"

	resp, err := http.PostForm(verifyLink,
		url.Values{"secret": {secretKey}, "response": {recaptcha}})
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	type JSONAPIResponse struct {
		Success bool `json:"success"`
	}
	var googleResponse JSONAPIResponse

	json.Unmarshal(body, &googleResponse)

	if googleResponse.Success == true {
		verifySuccess = "1"
		// fmt.Println("Google Recaptcha Success")
	} else {
		// fmt.Println("Google Recaptcha failure")
	}
	return verifySuccess
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
