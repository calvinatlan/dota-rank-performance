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

func checkPlayerExists(db *sql.DB, id string) bool {
	queryPlayer := `SELECT playerid FROM player WHERE playerid = $1`
	rows, _ := db.Query(queryPlayer, id)
	cols, _ := rows.Columns()
	if len(cols) == 0 {
		return false
	}
	return true
}

func updatePlayer(db *sql.DB, id string) {
	queryPlayer := `SELECT lastupdated FROM player WHERE playerid = $1`
	rows, _ := db.Query(queryPlayer, id)
	var lastUpdated time.Time
	rows.Next()
	rows.Scan(&lastUpdated)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var daysBack int
	if lastUpdated.Before(thirtyDaysAgo) {
		daysBack = 30
	} else {
		daysBack = int((time.Now().Sub(lastUpdated).Hours())/24) + 1
	}
	log.Print(lastUpdated)
	strDaysBack := strconv.Itoa(daysBack)
	log.Print("Days back: " + strDaysBack)
	/*
		resp, err := http.Get("https://api.opendota.com/api/players/" + id + "/matches?significant=1&date=" + strDaysBack)
		if err != nil {
			log.Fatalln("Could not call opendota api for this player")
			log.Fatalln(err)
		}
		type Match struct {
			StartTime int64 `json:"start_time"`
			MatchId   int `json:"match_id"`
		}
		var responseJson []Match
		byteValue, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(byteValue, &responseJson)
	*/
	queryUpdateLU := `UPDATE player SET lastupdated = $1 WHERE playerid = $2`
	db.Exec(queryUpdateLU, time.Now(), id)
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

		playerExists := checkPlayerExists(db, id)
		if !playerExists {
			insertPlayer(db, id)
		}
		updatePlayer(db, id)

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
