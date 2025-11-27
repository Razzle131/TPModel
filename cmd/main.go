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

func WithoutCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func main() {
	logger.SetupLogging(slog.LevelDebug)
	cfg := config.InitCfg()

	server := handler.NewServer()

	http.HandleFunc("POST /valve/open/{valve}", WithoutCORS(server.PostOpenValve))
	http.HandleFunc("POST /valve/close/{valve}", WithoutCORS(server.PostCloseValve))
	http.HandleFunc("POST /drp", WithoutCORS(server.PostToggleDRP))
	http.HandleFunc("GET /sensor/{sensor}", WithoutCORS(server.GetSensor))
	http.HandleFunc("GET /tank/{tank}", WithoutCORS(server.GetTank))
	http.HandleFunc("GET /valve/{valve}", WithoutCORS(server.GetValve))

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
