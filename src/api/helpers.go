package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func (app application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
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

func (app application) readJSON(w http.ResponseWriter, r *http.Request, target any) error {

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

func loadEnv(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := map[string]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return env, nil
}

func getEnvValue(env map[string]string, key string) (string, bool) {
	value, exists := env[key]
	return value, exists
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	timeout, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(timeout)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
