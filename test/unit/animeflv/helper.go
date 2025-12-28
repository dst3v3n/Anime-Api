// Package animeflv contiene funciones auxiliares para testing unitario del scraper de AnimeFlv.
// Incluye utilidades para cargar archivos de fixtures (HTML embebidos) utilizados en las pruebas.
package animeflv

import (
	"os"
	"path/filepath"
	"testing"
)

func loadFixtures(t *testing.T, ruta string) *os.File {
	t.Helper()
	path := filepath.Join("fixtures", ruta)
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("No se pudo abrir el archivo de fixtures %s: %v", ruta, err)
	}

	return file
}
