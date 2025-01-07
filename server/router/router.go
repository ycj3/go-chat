package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ycj3/go-chat/server/handlers"
	"github.com/ycj3/go-chat/server/websocket"
)

func NewRouter(userHandler *handlers.UserHandler, hub *websocket.Hub) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/login", userHandler.HandleLogin).Methods("POST")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	}).Methods("GET")
	r.HandleFunc("/online", func(w http.ResponseWriter, r *http.Request) {
		onlineInfo := map[string]interface{}{
			"count":   hub.GetOnlineCount(),
			"members": hub.GetOnlineMembers(),
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(onlineInfo)
	}).Methods("GET")
	return r
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
