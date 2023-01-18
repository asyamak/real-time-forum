package main

import (
	"flag"
	"real-time-forum/internal/app"
)

func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	app := app.New()
	app.Start(configPath, "sqlite")
}
