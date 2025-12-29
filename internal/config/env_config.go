// Package config carga y gestiona la configuración de la aplicación.
// Proporciona un singleton de configuración que carga variables de entorno
// al inicio de la aplicación y las valida. Incluye configuración de caché (Valkey)
// y logging (Zerolog) con valores por defecto si no están definidas las variables de entorno.
package config

import (
	"fmt"
	"os"
	"sync"
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
	EnableCache   bool
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

// NewConfigWithDefaults crea una nueva instancia de Config con valores por defecto.
// Retorna una Config preconfigurada sin necesidad de variables de entorno,
// útil para testing o ejecución local sin configuración externa.
func NewConfigWithDefaults() *Config {
	return &Config{
		AppName: "Anime-API",
		CacheConfig: CacheConfig{
			CacheHost:   "localhost",
			CachePort:   6379,
			CacheDB:     0,
			CacheTTL:    60,
			EnableCache: false,
		},
		LogConfig: LogConfig{
			LogAppName: "Anime-API",
			LogEnv:     "development",
		},
	}
}

// NewConfigFromEnv carga la configuración desde variables de entorno.
// Intenta cargar un archivo .env si existe, pero continúa si no se encuentra.
// Retorna un error si la validación de configuración falla.
func NewConfigFromEnv() (*Config, error) {
	return NewConfigFromEnvPath("")
}

// NewConfigFromEnvPath carga la configuración desde un archivo .env específico o variables de entorno.
// Si envPath está vacío, intenta cargar desde la ruta por defecto.
// Retorna error si el archivo existe pero no puede ser leído, o si la configuración no es válida.
func NewConfigFromEnvPath(envPath string) (*Config, error) {
	if err := loadEnvFileFromPath(envPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	cfg := &Config{
		AppName: getEnv("APP_NAME", "Anime-API"),
		CacheConfig: CacheConfig{
			CacheHost:     getEnv("CACHE_HOST", "localhost"),
			CachePort:     getEnvAsInt("CACHE_PORT", 6379),
			CacheUsername: getEnv("CACHE_USERNAME", ""),
			CachePassword: getEnv("CACHE_PASSWORD", ""),
			CacheDB:       getEnvAsInt("CACHE_DB", 0),
			CacheTTL:      getEnvAsInt("CACHE_TTL", 3600),
			EnableCache:   getEnvAsBool("CACHE_ENABLED", false),
		},
		LogConfig: LogConfig{
			LogAppName: getEnv("LOG_APP_NAME", "MyApp"),
			LogEnv:     getEnv("LOG_ENV", "development"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// WithCacheHost establece el host del servidor Valkey (caché distribuido).
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCacheHost(host string) *Config {
	c.CacheHost = host
	return c
}

// WithCachePort establece el puerto del servidor Valkey (por defecto 6379).
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCachePort(port int) *Config {
	c.CachePort = port
	return c
}

// WithCacheUsername establece el usuario para autenticación en Valkey.
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCacheUsername(username string) *Config {
	c.CacheUsername = username
	return c
}

// WithCachePassword establece la contraseña para autenticación en Valkey.
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCachePassword(password string) *Config {
	c.CachePassword = password
	return c
}

// WithCacheDB establece el número de base de datos Valkey (0-15, por defecto 0).
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCacheDB(db int) *Config {
	c.CacheDB = db
	return c
}

// WithCacheTTL establece el tiempo de vida de los valores en caché en minutos (por defecto 3600).
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCacheTTL(ttl int) *Config {
	c.CacheTTL = ttl
	return c
}

// WithCache establece si el caché está habilitado (deprecated, usa WithEnableCache).
// Retorna la Config para encadenamiento de métodos según el patrón Builder.
func (c *Config) WithCache(enabled bool) *Config {
	c.EnableCache = enabled
	return c
}

// InitConfig inicializa el singleton de configuración. Solo se puede ejecutar una vez.
// Las siguientes llamadas son ignoradas si la instancia ya fue inicializada.
// Retorna error si la configuración no valida o si el singleton ya fue inicializado con diferente Config.
func InitConfig(cfg *Config) error {
	var initErr error
	once.Do(func() {
		if err := cfg.validate(); err != nil {
			initErr = err
			return
		}
		instance = cfg
	})

	if initErr != nil {
		return initErr
	}

	if instance != cfg {
		return fmt.Errorf("config singleton already initialized")
	}

	return nil
}

// GetConfig retorna la instancia singleton de Config, inicializándola desde variables de entorno si no existe.
// Usa inicialización lazy con patrón once para garantizar que solo se carga una vez.
// Retorna error si no puede cargar la configuración desde el entorno.
func GetConfig() (*Config, error) {
	once.Do(func() {
		instance, loadErr = NewConfigFromEnv()
	})
	return instance, loadErr
}

// MustGetConfig retorna la instancia singleton de Config. Panica si hay error al cargar.
// Útil cuando se espera que la configuración siempre esté disponible.
func MustGetConfig() *Config {
	cfg, err := GetConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

// ResetConfig reinicia el singleton de configuración.
// Únicamente para uso en pruebas (testing). NO usar en producción.
func ResetConfig() {
	once = sync.Once{}
	instance = nil
	loadErr = nil
}
