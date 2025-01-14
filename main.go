package main

import (
	"flag"
	"fmt"
	"os"

	"go-chat/api"
	"go-chat/connect"
	"go-chat/logic"

	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	// Set log level to Debug
	logrus.SetLevel(logrus.DebugLevel)

	var module string
	flag.StringVar(&module, "module", "", "assign run module")
	flag.Parse()
	fmt.Println(fmt.Sprintf("start run %s module", module))

	switch module {
	case "logic":
		logic.New().Run()
	case "connect":
		logrus.Info("Starting server on address:", *addr)
		connect.NewConnect().Run()
	case "api":
		api.New().Run()
	default:
		fmt.Println("exiting,module param error!")
		return
	}
	fmt.Println(fmt.Sprintf("run %s module done!", module))

	// dsn := "root@tcp(127.0.0.1:3306)/Chat?charset=utf8mb4&parseTime=True&loc=Local"
	// logrus.Debug("Connecting to database with DSN:", dsn)
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	logrus.Fatal("Failed to connect to database:", err)
	// }
	// logrus.Info("Database connection established")

	// // Auto migrate the User schema
	// logrus.Info("Auto migrating User schema")
	// db.AutoMigrate(&models.User{})

	// hub := websocket.NewHub(db)
	// go hub.Run()
	// logrus.Info("WebSocket hub started")

	// userService := services.NewUserService(db)
	// userHandler := handlers.NewUserHandler(userService)
	// r := router.NewRouter(userHandler, hub)

	// logrus.Fatal(http.ListenAndServe(*addr, r))
}
