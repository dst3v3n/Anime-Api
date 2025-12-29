package config

import "github.com/dst3v3n/api-anime/internal/config"

// NewConfigWithDefaults crea una nueva instancia de Config con valores por defecto.
// Este es un wrapper que delega a la función equivalente en el package internal/config.
// Retorna una Config preconfigurada sin necesidad de variables de entorno,
// útil para testing o ejecución local sin configuración externa.
// Parámetros: ninguno
// Retorna: una instancia de Config con valores por defecto predefinidos
func NewConfigWithDefaults() *config.Config {
	return config.NewConfigWithDefaults()
}

// NewConfigFromEnvPath carga la configuración desde un archivo .env específico o variables de entorno.
// Este es un wrapper que delega a la función equivalente en el package internal/config.
// Comportamiento:
//   - Si path está vacío: intenta cargar desde la raíz del proyecto o rutas relativas comunes
//   - Si path tiene valor: carga directamente desde esa ruta específica
// Parámetros:
//   - path: ruta específica al archivo .env (si está vacío, usa búsqueda automática)
// Retorna: una instancia de Config cargada desde el archivo o variables de entorno, o error si falla
func NewConfigFromEnvPath(path string) (*config.Config, error) {
	return config.NewConfigFromEnvPath(path)
}

// InitConfig inicializa el singleton de configuración de forma thread-safe.
// Este es un wrapper que delega a la función equivalente en el package internal/config.
// Solo se puede ejecutar una vez; las siguientes llamadas retornan error si intentan usar una Config diferente.
// Parámetros:
//   - cfg: instancia de Config a inicializar como singleton global
// Retorna: nil si la inicialización es exitosa, error si ya fue inicializado con diferente Config
func InitConfig(cfg *config.Config) error {
	return config.InitConfig(cfg)
}

// MustGetConfig retorna la instancia singleton de Config, inicializándola si no existe.
// Este es un wrapper que delega a la función equivalente en el package internal/config.
// Panica si hay error al cargar la configuración desde el entorno.
// Útil cuando se espera que la configuración siempre esté disponible sin manejo de error.
// Parámetros: ninguno
// Retorna: la instancia singleton de Config (nunca retorna error)
func MustGetConfig() *config.Config {
	return config.MustGetConfig()
}

// ResetConfig reinicia el singleton de configuración a su estado inicial.
// Este es un wrapper que delega a la función equivalente en el package internal/config.
// ÚNICAMENTE para uso en pruebas (testing). NO usar en producción.
// Permite que los tests reinicialicen la configuración entre casos de prueba.
// Parámetros: ninguno
// Retorna: nada
func ResetConfig() {
	config.ResetConfig()
}
