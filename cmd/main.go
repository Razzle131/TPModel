package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/Razzle131/TPModel/internal/handler"
	"github.com/Razzle131/TPModel/pkg/config"
	"github.com/Razzle131/TPModel/pkg/logger"
)

func main() {
	logger.SetupLogging(slog.LevelDebug)
	cfg := config.InitCfg()

	server := handler.NewServer()

	http.HandleFunc("POST /valve/open/{valve}", server.PostOpenValve)
	http.HandleFunc("POST /valve/close/{valve}", server.PostCloseValve)
	http.HandleFunc("POST /drp", server.PostToggleDRP)
	http.HandleFunc("GET /sensor/{sensor}", server.GetSensor)
	http.HandleFunc("GET /tank/{tank}", server.GetTank)

	srv := http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	slog.Info(fmt.Sprintf("starting server on address: %s", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen and serve error: %s", err.Error())
	}
}
