// Package config contiene funciones de validación para la estructura de configuración.
package config

import "fmt"

// validate verifica que todos los parámetros de configuración sean válidos y cumplan con los requerimientos.
// Valida:
// - APP_NAME: debe estar definido y no estar vacío
// - CACHE_PORT: debe estar en el rango válido de puertos (1-65535)
// - CACHE_TTL: debe ser un número no negativo (en minutos)
// - LOG_ENV: debe ser uno de los valores permitidos (development, staging, production)
// Retorna un error descriptivo si alguna validación falla, o nil si todas las validaciones pasan.
func (c *Config) validate() error {
	if c.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}

	if c.CachePort < 1 || c.CachePort > 65535 {
		return fmt.Errorf("invalid CACHE_PORT: must be between 1-65535, got %d", c.CachePort)
	}

	if c.CacheTTL < 0 {
		return fmt.Errorf("CACHE_TTL must be positive, got %d", c.CacheTTL)
	}

	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.LogEnv] {
		return fmt.Errorf("invalid LOG_ENV: must be development, staging or production, got %s", c.LogEnv)
	}

	return nil
}
