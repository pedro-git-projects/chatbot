package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pedro-git-projects/chatbot-back/src/data"
	"github.com/sashabaranov/go-openai"
)

type Application struct {
	config       *Config
	logger       *log.Logger
	models       data.Models
	openaiClient openai.Client
	upgrader     websocket.Upgrader
	clients      map[*websocket.Conn]bool
	broadcast    chan []byte
}

func InitializeApp() (*Application, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &Application{
		config: cfg,
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}

	return app, nil
}

func (app *Application) Run() error {
	db, err := openDB(*app.config)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	app.models = data.NewModels(db)
	app.logger.Printf("Conexão com o banco de dados estabelecida\n")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go app.handleMessages()
	go app.handleConnections()

	app.logger.Printf("Inicializando servidor em modo de %s na porta %s", app.config.env, srv.Addr)
	err = srv.ListenAndServe()
	return err
}

func (app Application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app Application) readJSON(w http.ResponseWriter, r *http.Request, target any) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(target)

	if err != nil {
		syntaxError := &json.SyntaxError{}
		unmarshalTypeError := &json.UnmarshalTypeError{}
		invalidUnmarshalError := &json.InvalidUnmarshalError{}

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Corpo contém JSON malformado (no caractere %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Corpo contém JSON malformado")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Corpo contém JSON do tipo incorreto para o campo %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("Corpo contém JSON contém JSON do tipo incorrecto (no caractere %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("O corpo da requisição não deve estar vazio")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: campo desconhecido")
			return fmt.Errorf("Corpo contém campos desconhecidos %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("O corpo da requisição não deve ser maior que %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("O Corpo deve conter um único valor JSON")
	}
	return nil
}
