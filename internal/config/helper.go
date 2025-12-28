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
func getEnvAsInt(name string, defaultVal int) int {
	if value := os.Getenv(name); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultVal
}

// findProjectRoot busca el directorio raíz del proyecto Go recorriendo hacia arriba en la estructura de directorios
// hasta encontrar un archivo go.mod. Retorna la ruta absoluta del directorio raíz o un error si no se encuentra.

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
// Intenta encontrar el archivo .env en la raíz del proyecto (buscando go.mod).
// Si no encuentra la raíz, busca en rutas relativas comunes (.env, ../.env, ../../.env, ../../../.env).
// Retorna un error si no puede cargar el archivo, o nil si lo carga exitosamente.
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
