package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pedro-git-projects/chatbot-back/internal/data"
)

const version = "1.0.0"

type config struct {
	port      int
	env       string
	jwtSecret string
	db        struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	cfg := config{}

	flag.IntVar(&cfg.port, "port", 4000, "Porta do servidor da API")
	flag.StringVar(&cfg.env, "env", "desenvolvimento", "Ambiente (desenvolvimento|homologação|produção)")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "Número máximo de conexões abertas no PostgreSQL")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "Número máximo de conexões inativas no PostgreSQL")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "Tempo máximo de conexão inativa no PostgreSQL")

	flag.Parse()

	env, err := loadEnv(".env")
	if err != nil {
		log.Fatal("falha ao ler arquivo .env:", err)
		return
	}

	key := "DATABASE_URL"
	if value, exists := getEnvValue(env, key); exists {
		cfg.db.dsn = value
	} else {
		fmt.Printf("Chave %s não foi encontrada no arquivo .env\n", key)
	}

	key = "JWT_SECRET"
	if value, exists := getEnvValue(env, key); exists {
		cfg.jwtSecret = value
	} else {
		fmt.Printf("Chave %s não foi encontrada no arquivo .env\n", key)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("Conexão com o banco de dados estabelecida\n")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Inicializando servidor em modo de %s na porta %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
