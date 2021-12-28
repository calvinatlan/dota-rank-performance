package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	conn := createConn()
	setUpRouter(conn)
}

func createConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func insertPlayer(conn *pgx.Conn, id string) {
	insertPlayerSQL := "INSERT INTO player(playerId, createdDate) VALUES($1, $2)"
	_, err := conn.Exec(context.Background(), insertPlayerSQL, id, time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	log.Print("Inserted player id " + id)
}

func getPlayers(conn *pgx.Conn) string {
	getPlayersSQL := `SELECT playerId FROM player;`
	rows, err := conn.Query(context.Background(), getPlayersSQL)
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

func setUpRouter(conn *pgx.Conn) {
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

	router.GET("/get-players", getPlayersHandler(conn))

	router.POST("/refresh-games", postRefreshGamesHandler(conn))

	router.Run(":" + port)
}

func getPlayersHandler(conn *pgx.Conn) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		players := getPlayers(conn)
		c.String(http.StatusOK, players)
	}
	return fn
}

func postRefreshGamesHandler(conn *pgx.Conn) gin.HandlerFunc {
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

		insertPlayer(conn, id)

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
