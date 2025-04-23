// internal/config/config.go
package config

import (
	"os"
)

type AppConfig struct {
	ClusterName      string
	ServiceAddr      string
	ControlPlaneHost string
	HttpPort         string
}

func Load() *AppConfig {
	return &AppConfig{
		ClusterName:      getEnv("EMBRACER_CLUSTER_NAME", "dynamic_backend_service"),
		ServiceAddr:      getEnv("EMBRACER_ADDR", "embracer:8080"),
		ControlPlaneHost: getEnv("CONTROL_PLANE_HOST", "embracer-control-plane:2222"),
		HttpPort:         getEnv("HTTP_PORT", ":8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
