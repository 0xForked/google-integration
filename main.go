package main

import (
	"database/sql"
	"log"
	"slices"
	"time"

	"github.com/0xForked/goca/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

func createNewDBConn() *sql.DB {
	driver, source := "sqlite", "./db.sqlite3"
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var allowOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8000",
}

var allowHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"Authorization",
	"Cache-Control",
	"Origin",
	"Cookie",
}

func createNewEngine() *gin.Engine {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET, POST, PATCH, DELETE"},
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return slices.Contains(allowOrigins, origin)
		},
		MaxAge: 12 * time.Hour,
	}))
	return engine
}

func main() {
	db := createNewDBConn()
	engine := createNewEngine()
	server.Run(db, engine)
}
