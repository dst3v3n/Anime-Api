// Package config proporciona funciones para inicializar y obtener loggers configurados
// según los parámetros de la aplicación. Integra el sistema de logging con Zerolog
// permitiendo loggers globales que respeten la configuración del entorno (development, staging, production)
// y el nombre de la aplicación definidos en la Config.
package config

import (
	"github.com/dst3v3n/api-anime/pkg/logger"
	"github.com/rs/zerolog"
)

// Logging inicializa y retorna un logger de Zerolog configurado según los parámetros de Config.
// El logger se configura con el entorno y nombre de aplicación especificados.
func (c *Config) Logging() zerolog.Logger {
	return logger.InitLogger(c.LogEnv, c.LogAppName)
}

// GetLogger retorna el logger global de la aplicación basado en la configuración.
// Si hay error al obtener la configuración, retorna un logger con valores por defecto (development, Anime-API).
func GetLogger() zerolog.Logger {
	cfg, err := GetConfig()
	if err != nil {
		return logger.InitLogger("development", "Anime-API")
	}
	return cfg.Logging()
}
