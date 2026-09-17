package sqlstreams

import (
	"crypto/tls"
	"fmt"
	"time"
)

// PostgresConnectionConfig configures the pool created by NewPostgresPool.
type PostgresConnectionConfig struct {
	Port int // Default: 5432.
	// MaxConns sets the pool size when positive. Zero preserves pgx's default.
	MaxConns int
	// ConnectTimeout sets the connection timeout when positive.
	// Zero preserves the parsed pgx connection settings.
	ConnectTimeout time.Duration
	// TLSConfig replaces the parsed TLS configuration when non-nil.
	// It does not replace pgx's fallback settings. Use a pgx connection URL
	// to configure sslmode and fallback behavior explicitly.
	TLSConfig *tls.Config
}

// WithDefaults fills Port (5432) -- the one knob that's a protocol constant.
func (c *PostgresConnectionConfig) WithDefaults() *PostgresConnectionConfig {
	if c.Port == 0 {
		c.Port = 5432
	}
	return c
}

func (c *PostgresConnectionConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("Port must be > 0, got %d", c.Port)
	}
	if c.MaxConns < 0 {
		return fmt.Errorf("MaxConns must be >= 0, got %d", c.MaxConns)
	}
	if c.ConnectTimeout < 0 {
		return fmt.Errorf("ConnectTimeout must be >= 0, got %s", c.ConnectTimeout)
	}
	return nil
}
