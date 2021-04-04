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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var HOST string
var DATABASE string
var USER string
var PASSWORD string

func redicertShortLink(c *gin.Context) {
	shortLinkCode := strings.ReplaceAll((c.Request.URL.Path), "/", "")

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE),
	)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT link FROM shortlink WHERE code = $1", shortLinkCode)

	var linkcheck bool

	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			fmt.Println(err)
		}
		if link == "" {
			linkcheck = false
		} else {
			c.Redirect(302, link)
			linkcheck = true
		}
	}

	if linkcheck == false {
		c.HTML(404, "404.tmpl", nil)
	}
}

func shortLinkCreate(c *gin.Context) {

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
	token := c.PostForm("token")
	link := c.PostForm("link")
	admin := c.PostForm("admin")
	googleRecaptcha := c.PostForm("g-recaptcha-response")
	custom := c.PostForm("custom")
	customcode := c.PostForm("customcode")

	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") == true {
		if admin == "true" {
			if token == tokenValue {

				var code string
				code = randomString(3)

				if custom == "true" {
					code = customcode
				}

				if createShortLink(code, link, c.ClientIP()) == true {
					r = Result{true, "Succeed to create Link!", "https://yiy.tw/" + code}
					c.JSON(200, r)
				} else {
					r = Result{false, "Failure", ""}
					c.JSON(400, r)
				}

			} else {
				r = Result{false, "Unauthorized", ""}
				c.JSON(403, r)
			}
		} else {
			if verifyRecaptcha(googleRecaptcha) == true {
				var code string
				code = randomString(3)
				var checkforLoop bool

				for checkforLoop == false {
					if checkCodeAvailable(code) == false {
						if createShortLink(code, link, c.ClientIP()) == true {
							r = Result{true, "Succeed to create Link!", "https://yiy.tw/" + code}
							c.JSON(200, r)
						} else {
							r = Result{false, "Failure", ""}
							c.JSON(400, r)
						}
						checkforLoop = true
					} else {
						code = randomString(3)
					}
				}

			} else {
				r = Result{false, "Google Recaptcha Failure!", ""}
				c.JSON(403, r)
			}
		}
	} else {
		r = Result{false, "Please type a Vaild URL.", ""}
		c.JSON(400, r)
	}
}

func checkCodeAvailable(code string) bool {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE),
	)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT link FROM shortlink WHERE code = $1", "steveyi")

	var check bool

	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			fmt.Println(err)
		}
		if link == "" {
			check = false
		} else {
			check = true
		}
	}
	return check
}

func createShortLink(code string, link string, IP string) bool {

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE),
	)
	if err != nil {
		panic(err)
	}

	sqlStatement := `
	INSERT INTO shortlink(code, link, ipaddress) values($1,$2,$3)`
	_, err = db.Exec(sqlStatement, code, link, IP)

	if err != nil {
		panic(err)
	}

	var makeshortlink bool

	if err != nil {
		fmt.Println(err)
	} else {
		makeshortlink = true
	}
	return makeshortlink
}

func randomString(length int) string {
	rand.Seed(time.Now().Unix())

	var output strings.Builder

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJULMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return (output.String())
}

func verifyRecaptcha(recaptcha string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	verifyLink := "https://www.recaptcha.net/recaptcha/api/siteverify"
	secretKey := os.Getenv("Google_Recaptcha_SecretKey")
	verifySuccess := false

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
		verifySuccess = true
	}
	return verifySuccess
}

func indexPage(c *gin.Context) {
	c.HTML(200, "index.tmpl", nil)
}

func pageNotAvailable(c *gin.Context) {
	c.HTML(404, "404.tmpl", nil)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	HOST = os.Getenv("SQL_HOST")
	DATABASE = os.Getenv("SQL_DATABASE")
	USER = os.Getenv("SQL_USER")
	PASSWORD = os.Getenv("SQL_PASSWORD")

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.LoadHTMLGlob("static/*")

	router.GET("/", indexPage)
	router.GET("/:ShortLinkCode", redicertShortLink)
	router.POST("/api/v1/create", shortLinkCreate)

	router.NoRoute(pageNotAvailable)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	router.Run(":32156")
}
