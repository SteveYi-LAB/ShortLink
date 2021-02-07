package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func webServer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Method:", r.Method)
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path

	if r.Method == "GET" {
		if p == "./" || p == "./index.html" {
			http.ServeFile(w, r, "./index.html")
		} else {
			shortLink := strings.ReplaceAll(p, "./", "")
			fmt.Println(shortLink)
		}
	}

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		fmt.Println("Method:", r.Method)
		fmt.Println(r.URL.Path)

		if r.URL.Path == "/api/create" {
			tokenValue := "123456"

			r.ParseForm()
			token := "NULL"
			link := "NULL"
			fmt.Println(link)

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
				fmt.Println(code)

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
