package server

import (
	api "signupin-api/internal/app/api"
)

// NewServer Return new server instance
func NewServer() {
	api.CreateAPIApp()
}
