package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func shortLinkCreate(ctx *gin.Context) {

	type Result struct {
		Success bool
		Message string
		Link    string
	}

	var r Result

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenValue := os.Getenv("token")

	token := ctx.PostForm("token")
	link := ctx.PostForm("link")
	admin := ctx.PostForm("admin")
	googleRecaptcha := ctx.PostForm("g-recaptcha-response")
	custom := ctx.PostForm("custom")
	customcode := ctx.PostForm("customcode")

	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") == true {
		if admin == "true" {
			if token == tokenValue {

				var code string
				code = randomString(3)

				if custom == "true" {
					code = customcode
				}

				if createShortLink(code, link, ctx.ClientIP()) == 1 {
					r = Result{true, "Succeed to create Link!", "https://yiy.tw/" + code}
				} else {
					r = Result{false, "Failure", ""}
				}

			} else {
				r = Result{false, "Unauthorized", ""}
			}
		} else {
			if verifyRecaptcha(googleRecaptcha) == "1" {
				var code string
				code = randomString(3)
				var checkforLoop int

				for checkforLoop == 0 {
					if checkCodeAvailable(code) == 0 {
						if createShortLink(code, link, ctx.ClientIP()) == 1 {
							r = Result{true, "Succeed to create Link!", "https://yiy.tw/" + code}
						} else {
							r = Result{false, "Failure", ""}
						}
						checkforLoop = 1
					}
				}

			} else {
				r = Result{false, "Google Recaptcha Failure!", ""}
			}
		}
	} else {
		r = Result{false, "Please type a Vaild URL.", ""}
	}
	ctx.JSON(200, r)
}

func redicertShortLink(ctx *gin.Context) {
	if ctx.Request.URL.Path == "/" {
		ctx.HTML(302, "index.html", nil)
	} else {
		shortLinkCode := strings.ReplaceAll((ctx.Request.URL.Path), "/", "")
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
				ctx.Redirect(302, link)
				linkcheck = 1
			}
		}

		if linkcheck == 0 {
			ctx.HTML(404, "404.html", nil)
		}
	}
}

func pageNotAvailable(ctx *gin.Context) {
	ctx.HTML(404, "404.html", nil)
}

func createShortLink(code string, link string, IP string) int {

	db, err := sql.Open("sqlite3", "db")
	if err != nil {
		fmt.Println(err)
	}

	stmt, err := db.Prepare("INSERT INTO shortlink(code, link, ipAddress) values(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(code, link, IP)
	if err != nil {
		fmt.Println(err)
	}
	id, err := res.LastInsertId()
	var makeshortlink int

	if err != nil {
		fmt.Println(err)
		fmt.Println(id)
	} else {
		makeshortlink = 1
	}
	return makeshortlink
}

func checkCodeAvailable(code string) int {

	db, err := sql.Open("sqlite3", "db")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := db.Query("SELECT link FROM shortlink WHERE code = ?", code)

	var check int

	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			fmt.Println(err)
		}
		if link == "" {
			check = 0
		} else {
			check = 1
		}
	}
	return check
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
	}
	return verifySuccess
}

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.LoadHTMLGlob("static/*")

	router.GET("/:ShortLinkCode", redicertShortLink)
	router.POST("/api/create", shortLinkCreate)

	router.NoRoute(pageNotAvailable)

	router.Run(":32156")
}
