package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	db := createConn()
	setUpRouter(db)
}

func createConn() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func insertPlayer(db *sql.DB, id string) {
	insertPlayerSQL := "INSERT INTO player(playerid, createddate) VALUES($1, $2)"
	_, err := db.Exec(insertPlayerSQL, id, time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	log.Print("Inserted player id " + id)
}

func getPlayers(db *sql.DB) string {
	getPlayersSQL := `SELECT playerid FROM player;`
	rows, err := db.Query(getPlayersSQL)
	if err != nil {
		log.Fatalln(err)
	}
	var id int
	var ids []string
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, strconv.Itoa(id))
	}
	return strings.Join(ids, " ")
}

func setUpRouter(db *sql.DB) {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/get-players", getPlayersHandler(db))

	router.POST("/refresh-games", postRefreshGamesHandler(db))

	router.Run(":" + port)
}

func getPlayersHandler(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		players := getPlayers(db)
		c.String(http.StatusOK, players)
	}
	return fn
}

func postRefreshGamesHandler(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		type RefreshGameBody struct {
			PlayerId string `json:"playerId"`
		}
		var requestBody RefreshGameBody
		err := c.BindJSON(&requestBody)
		if err != nil {
			log.Fatalln(err)
		}
		id := requestBody.PlayerId

		insertPlayer(db, id)

		/*resp, err := http.Get("https://api.opendota.com/api/players/" + id + "/peers")
		if err != nil {
			log.Fatalln(err)
			c.String(http.StatusInternalServerError, "")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			c.String(http.StatusInternalServerError, "")
		}
		//c.String(http.StatusOK, string(body))
		*/
		c.String(http.StatusOK, string("Refreshed player "+id))

	}
	return fn
}
