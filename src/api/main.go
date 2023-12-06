package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	cfg := config{}

	flag.IntVar(&cfg.port, "port", 4000, "Porta do servidor da API")
	flag.StringVar(&cfg.env, "env", "desenvolvimento", "Ambiente (desenvolvimento|homologação|produção)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Inicializando servidor em modo de %s na porta %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
