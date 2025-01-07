package main

import (
	"flag"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/ycj3/go-chat/server/handlers"
	"github.com/ycj3/go-chat/server/models"
	"github.com/ycj3/go-chat/server/router"
	"github.com/ycj3/go-chat/server/services"
	"github.com/ycj3/go-chat/server/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	// Set log level to Debug
	logrus.SetLevel(logrus.DebugLevel)

	flag.Parse()
	logrus.Info("Starting server on address:", *addr)

	dsn := "root@tcp(127.0.0.1:3306)/Chat?charset=utf8mb4&parseTime=True&loc=Local"
	logrus.Debug("Connecting to database with DSN:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Failed to connect to database:", err)
	}
	logrus.Info("Database connection established")

	// Auto migrate the User schema
	logrus.Info("Auto migrating User schema")
	db.AutoMigrate(&models.User{})

	hub := websocket.NewHub(db)
	go hub.Run()
	logrus.Info("WebSocket hub started")

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)
	r := router.NewRouter(userHandler, hub)

	logrus.Fatal(http.ListenAndServe(*addr, r))
}
