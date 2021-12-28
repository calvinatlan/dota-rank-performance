package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("PORT")

	id := "29449073"

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/refresh-games", func(c *gin.Context) {
		resp, err := http.Get("https://api.opendota.com/api/players/" + id + "/peers")
		if err != nil {
			log.Fatalln(err)
			c.String(http.StatusInternalServerError, "")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			c.String(http.StatusInternalServerError, "")
		}
		c.String(http.StatusOK, string(body))
	})

	router.Run(":" + port)
}
