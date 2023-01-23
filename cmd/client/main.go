package main

import (
	"flag"
	"html/template"
	"net/http"
	"real-time-forum/internal/config"
	"real-time-forum/pkg/logger"
)

func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	log := logger.NewLogger("[Forum]")

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Error(err.Error())
	}

	indexTemp, err := template.ParseFiles("./web/public/index.html")
	if err != nil {
		log.Error("%e", err)
	}

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./web/src"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := indexTemp.Execute(w, cfg.ServerAddress()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	log.Info("Client server is started at %v", cfg.Client.Port)
	if err := http.ListenAndServe(":"+cfg.Client.Port, nil); err != nil {
		log.Error("client server: error shile starting server: %e", err.Error())
	}
}
