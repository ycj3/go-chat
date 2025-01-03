package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/ycj3/go-chat/server/models"
	"github.com/ycj3/go-chat/server/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("serveHome called with URL:", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

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

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("WebSocket connection request received")
		websocket.ServeWs(hub, w, r)
	})
	http.HandleFunc("/online", func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("Online users request received")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		onlineInfo := map[string]interface{}{
			"count":   hub.GetOnlineCount(),
			"members": hub.GetOnlineMembers(),
		}
		json.NewEncoder(w).Encode(onlineInfo)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		logrus.Fatal("ListenAndServe: ", err)
	}
}
