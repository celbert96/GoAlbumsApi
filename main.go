package main

import (
	"database/sql"
	"middleware"
	"models"
	"os"
	"routes"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	//initialize db
	cfg := mysql.Config{
		User:   os.Getenv("MYSQLTEST_DB_USER"),
		Passwd: os.Getenv("MYSQLTEST_DB_PASS"),
		Net:    "tcp",
		Addr:   "10.0.0.32:3306",
		DBName: "mysqltest",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	env := &models.Env{DB: db}

	router := gin.Default()
	router.Use(middleware.EnvMiddleware(*env))

	/* Public Routes */
	pubv1 := router.Group("/v1")
	routes.AddAuthRoutes(pubv1)

	/* Private Routes */
	privv1 := router.Group("/v1")
	privv1.Use(middleware.TokenAuthMiddleware())
	routes.AddAlbumRoutes(privv1)

	router.Run(":8080")
}
