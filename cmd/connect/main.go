package main

import (
	"go-chat/internal/connect"
	"go-chat/pkg/utils"
)

func main() {
	utils.InitLogrus()

	connect.NewConnect().Run()
}
