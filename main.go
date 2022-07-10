package main

import (
	"database/sql"
	"gin-learning/middleware"
	"gin-learning/models"
	"gin-learning/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	//initialize db
	cfg := mysql.Config{
		User:   os.Getenv("MYSQLTEST_DB_USER"),
		Passwd: os.Getenv("MYSQLTEST_DB_PASS"),
		Net:    "tcp",
		Addr:   "192.168.1.17:3306",
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

	/* Protected routes which support cookie auth scheme */
	webprotectedv1 := router.Group("/web/v1")
	webprotectedv1.Use(middleware.CookieTokenAuth())
	routes.AddAlbumRoutes(webprotectedv1)

	/* Protected routes which support Bearer token auth scheme */
	apiprotectedv1 := router.Group("api/v1")
	apiprotectedv1.Use(middleware.BearerTokenAuth())
	routes.AddAlbumRoutes(apiprotectedv1)
	router.Run(":8080")
}
