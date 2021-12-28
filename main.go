package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	db := createDB()
	createTable(db)
	setUpRouter(db)
}

func createDB() *sql.DB {
	if _, err := os.Stat("dota-rank-performance.db"); os.IsNotExist(err) {
		file, err := os.Create("dota-rank-performance.db")
		if err != nil {
			log.Fatalln(err.Error())
		}
		file.Close()
	}

	db, _ := sql.Open("sqlite3", "dota-rank-performance.db")
	return db
}

func createTable(db *sql.DB) {
	createPlayerTableSQL := `
	CREATE TABLE IF NOT EXISTS player (
		"idPlayer" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"createdDate" DATE
	);`

	stmt, err := db.Prepare(createPlayerTableSQL)
	if err != nil {
		log.Fatalln(err)
	}
	stmt.Exec()
}

func insertPlayer(db *sql.DB, id string) {
	//insertPlayerSQL := `INSERT INTO player(idPlayer, createdDate) VALUES (?, ?);`
	insertPlayerSQL := `INSERT INTO player(idPlayer, createdDate) VALUES (?, ?);`
	stmt, err := db.Prepare(insertPlayerSQL)
	if err != nil {
		log.Fatalln(err)
	}
	stmt.Exec(id, time.Now())
	log.Print("Inserted player id " + id)
}

func getPlayers(db *sql.DB) string {
	getPlayersSQL := `SELECT idPlayer FROM player;`
	rows, err := db.Query(getPlayersSQL)
	if err != nil {
		log.Fatalln(err)
	}
	var id int
	var ids []string
	for rows.Next() {
		rows.Scan(&id)
		log.Print(id)
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
		log.Print(players)
		c.String(http.StatusOK, players)
	}
	return gin.HandlerFunc(fn)
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
	return gin.HandlerFunc(fn)
}
