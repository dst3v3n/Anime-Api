// Package config contiene funciones auxiliares para cargar y gestionar variables de entorno.
// Proporciona métodos para obtener variables de entorno con valores por defecto,
// localizar el archivo .env en la estructura del proyecto y cargar las configuraciones.
package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// getEnv obtiene el valor de una variable de entorno con un valor por defecto.
// Si la variable existe y no está vacía, retorna su valor; de lo contrario, retorna el defaultVal.
func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultVal
}

// getEnvAsInt obtiene el valor de una variable de entorno como entero con un valor por defecto.
// Si la variable existe y puede convertirse a entero válido, retorna ese valor; de lo contrario, retorna defaultVal.
// Parámetros:
//   - name: nombre de la variable de entorno a buscar
//   - defaultVal: valor por defecto si la conversión falla o la variable no existe
// Retorna: el valor convertido a entero o el valor por defecto
func getEnvAsInt(name string, defaultVal int) int {
	if value := os.Getenv(name); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultVal
}

// getEnvAsBool obtiene el valor de una variable de entorno como booleano con un valor por defecto.
// Interpreta valores como "true", "1", "yes" como verdadero y sus opuestos como falso.
// Si la variable existe y puede convertirse a booleano válido, retorna ese valor; de lo contrario, retorna defaultVal.
// Parámetros:
//   - name: nombre de la variable de entorno a buscar
//   - defaultVal: valor por defecto si la conversión falla o la variable no existe
// Retorna: el valor convertido a booleano o el valor por defecto
func getEnvAsBool(name string, defaultVal bool) bool {
	if value := os.Getenv(name); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}

	return defaultVal
}

// findProjectRoot busca el directorio raíz del proyecto Go recorriendo hacia arriba en la estructura de directorios
// hasta encontrar un archivo go.mod. Comienza desde el directorio de trabajo actual y sube recursivamente
// hacia los directorios padres hasta encontrar el archivo go.mod o llegar a la raíz del sistema de archivos.
// Parámetros: ninguno
// Retorna: la ruta absoluta del directorio raíz (donde está go.mod) o error si no se encuentra
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// loadEnvFile carga las variables de entorno desde un archivo .env.
// Estrategia de búsqueda:
//   1. Intenta encontrar la raíz del proyecto buscando go.mod usando findProjectRoot()
//   2. Si encuentra la raíz, carga .env desde ese directorio raíz
//   3. Si no encuentra la raíz, busca en rutas relativas comunes respecto al directorio actual:
//      - .env (mismo directorio)
//      - ../.env (un nivel arriba)
//      - ../../.env (dos niveles arriba)
//      - ../../../.env (tres niveles arriba)
// Parámetros: ninguno
// Retorna: nil si carga exitosamente el archivo, error si no puede cargar o si no lo encuentra en ninguna ruta
func loadEnvFile() error {
	root, err := findProjectRoot()
	if err != nil {
		paths := []string{".env", "../.env", "../../.env", "../../../.env"}
		for _, path := range paths {
			if err := godotenv.Load(path); err == nil {
				return nil
			}
		}
		return os.ErrNotExist
	}

	envPath := filepath.Join(root, ".env")
	return godotenv.Load(envPath)
}

// loadEnvFileFromPath carga variables de entorno desde una ruta específica o mediante búsqueda automática.
// Comportamiento:
//   - Si envPath está vacío: delega a loadEnvFile() para una búsqueda automática en la estructura del proyecto
//   - Si envPath tiene valor: intenta cargar directamente desde esa ruta específica del archivo .env
// Parámetros:
//   - envPath: ruta específica al archivo .env (si está vacío, usa búsqueda automática)
// Retorna: nil si carga exitosamente, error si no puede cargar el archivo desde la ruta especificada
func loadEnvFileFromPath(envPath string) error {
	if envPath == "" {
		return loadEnvFile() // Usa la búsqueda automática
	}
	return godotenv.Load(envPath)
}
