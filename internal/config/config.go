package config

import (
	"flag"
	"fmt"
)

// Config provides Storage, Address and DSN.
type Config struct {
	Storage string
	Addr    string
	DSN     string
}

// Supported storage types.
const (
	StorageMemory   = "memory"
	StoragePostgres = "postgres"
)

// Load reads and validates flags for main application.
func Load() (*Config, error) {
	storage := flag.String("storage", StorageMemory, "storage type: memory|postgres")
	addr := flag.String("addr", ":8080", "Address")
	dsn := flag.String("dsn", "", "Database connection string")
	flag.Parse()

	if *storage != StorageMemory && *storage != StoragePostgres {
		return nil, fmt.Errorf("unknown storage %q", *storage)
	}

	if *storage == StoragePostgres && *dsn == "" {
		return nil, fmt.Errorf("DSN is required for postgres storage")
	}
	return &Config{Storage: *storage, Addr: *addr, DSN: *dsn}, nil
}
