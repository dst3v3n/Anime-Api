// Package config carga y gestiona la configuración de la aplicación.
// Proporciona un singleton de configuración que carga variables de entorno
// al inicio de la aplicación y las valida. Incluye configuración de caché (Valkey)
// y logging (Zerolog) con valores por defecto si no están definidas las variables de entorno.
package config

import (
	"os"
	"sync"

	"github.com/dst3v3n/api-anime/pkg/logger"
	"github.com/rs/zerolog"
)

// Config contiene la configuración principal de la aplicación.
type Config struct {
	AppName string
	CacheConfig
	LogConfig
}

// CacheConfig contiene la configuración para la conexión a Valkey (caché distribuido).
type CacheConfig struct {
	CacheHost     string // Host del servidor Valkey
	CachePort     int    // Puerto del servidor Valkey
	CacheUsername string // Usuario para autenticación en Valkey
	CachePassword string // Contraseña para autenticación en Valkey
	CacheDB       int    // Número de base de datos Valkey
	CacheTTL      int    // Tiempo de vida de los valores en caché (en minutos)
}

// LogConfig contiene la configuración para el sistema de logging.
type LogConfig struct {
	LogAppName string // Nombre de la aplicación para los logs
	LogEnv     string // Entorno de ejecución (development, staging, production)
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

// NewConfig crea una nueva instancia de Config sin inicializar.
func NewConfig() *Config {
	return &Config{}
}

// getEnvieroment carga la configuración desde variables de entorno y archivo .env.
// Retorna una estructura Config con todos los parámetros cargados, o error si ocurre un problema.
func (c *Config) getEnvieroment() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	cfg := &Config{
		AppName: getEnv("APP_NAME", ""),

		CacheConfig: CacheConfig{
			CacheHost:     getEnv("CACHE_HOST", "localhost"),
			CachePort:     getEnvAsInt("CACHE_PORT", 6379),
			CacheUsername: getEnv("CACHE_USERNAME", ""),
			CachePassword: getEnv("CACHE_PASSWORD", ""),
			CacheDB:       getEnvAsInt("CACHE_DB", 0),
			CacheTTL:      getEnvAsInt("CACHE_TTL", 3600),
		},

		LogConfig: LogConfig{
			LogAppName: getEnv("LOG_APP_NAME", "MyApp"),
			LogEnv:     getEnv("LOG_ENV", "development"),
		},
	}

	return cfg, nil
}

// Logging retorna un logger configurado según el ambiente (development, staging, production).
// Utiliza la configuración global para inicializar zerolog con los parámetros apropiados.
func (c *Config) Logging() zerolog.Logger {
	cfg, err := GetConfig()
	if err != nil {
		return zerolog.Logger{}
	}
	return logger.InitLogger(cfg.LogEnv, cfg.LogAppName)
}

// GetConfig retorna la instancia singleton de configuración.
// Utiliza un patrón de sincronización (sync.Once) para garantizar que
// la configuración se cargue una única vez desde variables de entorno y se valide.
// Retorna error si la carga o validación falla.
func GetConfig() (*Config, error) {
	once.Do(func() {
		cfg := &Config{}
		instance, loadErr = cfg.getEnvieroment()
		if loadErr != nil {
			return
		}
		loadErr = instance.validate()
	})
	return instance, loadErr
}
