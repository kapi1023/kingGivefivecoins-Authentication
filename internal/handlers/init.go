package handlers

import (
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/config"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/storage"
)

var db *storage.PostgresStorage
var secret string
var cfg *config.Config

func InitializeHandlers(database *storage.PostgresStorage, jwtSecret string) {
	db = database
	secret = jwtSecret
	cfg = config.Load()
}
