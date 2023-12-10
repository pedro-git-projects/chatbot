package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/pedro-git-projects/chatbot-back/src/app"
)

func main() {
	app, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("A aplicação falhou com erro %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("A aplicação falhou com erro ao executar: %v", err)
	}
}
