package config

import (
	"fmt"
	"os"

	"github.com/calvinbui/teleport-sre-challenge/services/logger"
	"github.com/jessevdk/go-flags"
)

type Config struct {
	GRPCAddress string `long:"grpc-address" default:"0.0.0.0:8080" description:"address to run the grpc server" env:"GRPC_ADDRESS"`

	HttpAddress string `long:"http-address" default:"0.0.0.0:8081" description:"address to run the http server" env:"HTTP_ADDRESS"`

	FilePath string `long:"file-path" default:"/mnt/replicas.json" description:"where to store the state file" env:"FILE_PATH"`

	CertFile   string `long:"cert-file" default:"/certificates/server.crt" description:"Server public key"`
	KeyFile    string `long:"key-file" default:"/certificates/server.key" description:"Server private key"`
	CACertFile string `long:"ca-cert-file" default:"/certificates/ca.crt" description:"Root CA public key"`

	WatchNamespace string `long:"namespace" default:"" env:"NAMESPACE" description:"Namespace to watch. Leave empty to watch all namespaces"`

	LogLevel string `long:"log-level" default:"Info" env:"LOG_LEVEL" description:"Log level. Options are Info and Debug"`
}

// New gets the configuration using flags, env, etc.
func New() (Config, error) {
	var conf Config
	var parser = flags.NewParser(&conf, 0)
	if _, err := parser.Parse(); err != nil {
		return Config{}, fmt.Errorf("Error getting config: %s", err)
	}

	// set up log level now
	if err := logger.SetLevel(conf.LogLevel); err != nil {
		return Config{}, fmt.Errorf("Error setting log level: %w", err)
	}

	return conf, nil
}

// WriteHelp prints the help message
func WriteHelp() {
	var conf Config
	var parser = flags.NewParser(&conf, 0)
	parser.WriteHelp(os.Stdout)
}
