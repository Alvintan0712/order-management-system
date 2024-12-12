package main

import (
	"log"

	"example.com/oms/common"
	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	serviceName = common.EnvString("SERVICE_NAME", "gateway")
	serviceHost = common.EnvString("SERVICE_HOST", "127.0.0.1")
	servicePort = common.EnvString("SERVICE_PORT", "8080")

	debug = common.EnvString("DEBUG", "false") == "true"
)

func main() {
	config := &Config{
		consulAddr: consulAddr,
		Name:       serviceName,
		Host:       serviceHost,
		Port:       servicePort,
	}

	app := NewApp(config)
	defer app.Close()

	if err := app.Listen(); err != nil {
		log.Fatal("Failed to start http server:", err)
	}
}
