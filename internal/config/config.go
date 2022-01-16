package config

import (
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	// GpgTty string
	Prefix          string
	GpgRecipientsId []string
	// Extensions string
	// XSelection string
	// ClipTime string
	// GeneratedLength string
	// CharacterSet string
	// CharacterSetNoSymbols string
}

func New() config {
	cfg := config{}
	if cfg.Prefix = os.Getenv("PREFIX"); cfg.Prefix == "" {
		var home string
		if home = os.Getenv("HOME"); home == "" {
			panic("no HOME env set")
		}
		cfg.Prefix = filepath.Join(home, ".password-store")
	}
	if store := os.Getenv("PASSWORD_STORE_KEY"); store != "" {
		cfg.GpgRecipientsId = strings.Split(store, " ")
	}
	return cfg
}
