package main

import "github.com/VadimGossip/grpcAuditLog/internal/app"

var configDir = "config"

func main() {
	app.Run(configDir)
}
