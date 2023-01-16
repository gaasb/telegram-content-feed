package utils

import (
	"github.com/joho/godotenv"
	"log"
)

const (
	ENV_ERROR     = "Missing .env file."
	ENV_FILE      = ".env"
	TOKEN_ERROR   = "Missing '" + TOKEN + "' value in .env file!."
	DB_CONN_ERROR = "Missing '" + DB_CONNECT_URI + "' value in .env file!."
)

var (
	dotenv map[string]string
)

type Utilities[T any, E any] struct {
	DatabaseClient T
	BotClient      E
}

type Client[T any] struct {
	Value *T
	Setter
}

type Setter interface {
	Set(envValue string)
}

func NewUtilities[T any, E any](database *Client[T], telegram *Client[E]) *Utilities[T, E] {
	database.Set(dotenv[DB_CONNECT_URI])
	telegram.Set(dotenv[TOKEN])
	return &Utilities[T, E]{
		DatabaseClient: *database.Value,
		BotClient:      *telegram.Value,
	}
}
func GetEnv(value string) string {
	return dotenv[value]
}

func init() {
	env, err := godotenv.Read(ENV_FILE)
	if err != nil {
		log.Fatal(ENV_ERROR)
		return
	}
	dotenv = env
}
