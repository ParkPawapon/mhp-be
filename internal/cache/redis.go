package cache

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/ParkPawapon/mhp-be/internal/config"
)

func New(cfg config.RedisConfig) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.TLSEnabled {
		tlsConfig, err := buildTLSConfig(cfg)
		if err != nil {
			return nil, err
		}
		options.TLSConfig = tlsConfig
	}

	client := redis.NewClient(options)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func buildTLSConfig(cfg config.RedisConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.TLSInsecureSkipVerify,
	}

	if strings.TrimSpace(cfg.TLSServerName) != "" {
		tlsConfig.ServerName = strings.TrimSpace(cfg.TLSServerName)
	}

	if strings.TrimSpace(cfg.TLSCAFile) != "" {
		certPEM, err := os.ReadFile(cfg.TLSCAFile)
		if err != nil {
			return nil, err
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(certPEM) {
			return nil, fmt.Errorf("invalid redis tls ca file: %s", cfg.TLSCAFile)
		}
		tlsConfig.RootCAs = certPool
	}

	return tlsConfig, nil
}
