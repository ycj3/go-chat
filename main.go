package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/ycj3/go-chat/server/models"
	"github.com/ycj3/go-chat/server/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
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
	flag.Parse()

	dsn := "root@tcp(127.0.0.1:3306)/Chat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Auto migrate the User schema
	db.AutoMigrate(&models.User{})

	hub := websocket.NewHub(db)
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})
	http.HandleFunc("/online", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		onlineInfo := map[string]interface{}{
			"count":   hub.GetOnlineCount(),
			"members": hub.GetOnlineMembers(),
		}
		json.NewEncoder(w).Encode(onlineInfo)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
