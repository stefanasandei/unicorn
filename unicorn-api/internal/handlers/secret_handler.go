package handlers

import (
	"unicorn-api/internal/config"
	"unicorn-api/internal/stores"
)

func NewSecretsHandler(store *stores.SecretStore, iam stores.IAMStore, cfg *config.Config) *SecretsHandler {
	return &SecretsHandler{
		Store:    store,
		IAMStore: iam,
		Config:   cfg,
	}
}
