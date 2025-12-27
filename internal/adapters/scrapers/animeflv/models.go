// Package animeflv - models.go
// Este archivo define las estructuras de datos internas utilizadas específicamente
// por el scraper de AnimeFlv. Incluye:
// - Config: Configuración de URLs del sitio
// - ParseResult: Estructura temporal para almacenar datos durante el parsing de información de anime
// - ParseEpisodeLinksResult: Estructura temporal para almacenar enlaces de episodios
// - VideoServer y Videos: Estructuras para deserializar JSON embebido en scripts del sitio
package animeflv

import "github.com/dst3v3n/api-anime/internal/domain/dto"

// Config contiene la configuración de URLs del sitio AnimeFlv.
type Config struct {
	BaseURL       string // URL base del sitio (ej: "https://www3.animeflv.net")
	SearchURL     string // URL del endpoint de búsqueda
	AnimeInfoURL  string // URL del endpoint de información de anime
	VerEpisodeURL string // URL base para ver episodios
}

// ParseResult almacena temporalmente los datos extraídos durante el parsing de información de anime.
// Se utiliza como estructura intermedia antes de convertir a AnimeInfoResponse.
type ParseResult struct {
	title        string             // Título del anime
	category     string             // Categoría/Tipo del anime
	sipnopsis    string             // Sinopsis del anime
	status       string             // Estado de emisión del anime
	image        string             // URL de la imagen/carátula
	punctuacion  float64            // Calificación del anime
	animeRelated []dto.AnimeRelated // Animes relacionados
	genres       []string           // Géneros del anime
	episodes     []int              // Lista de episodios disponibles
	nextEpisode  string             // Fecha del próximo episodio
}

// ParseEpisodeLinksResult almacena temporalmente los enlaces extraídos de un episodio.
// Se utiliza como estructura intermedia antes de convertir a LinkResponse.
type ParseEpisodeLinksResult struct {
	ID      string           // Identificador del anime
	Title   string           // Título del anime
	Episode uint             // Número del episodio
	links   []dto.LinkSource // Enlaces de reproducción disponibles
}

// VideoServer representa un servidor de video individual con sus propiedades.
// Se utiliza para deserializar el JSON embebido en los scripts de AnimeFlv.
type VideoServer struct {
	Server      string `json:"server"`         // Nombre del servidor de video
	Title       string `json:"title"`          // Título/descripción del servidor
	Ads         int    `json:"ads"`            // Número de anuncios
	URL         string `json:"url,omitempty"`  // URL del enlace
	AllowMobile bool   `json:"allow_mobile"`   // Si permite acceso desde móviles
	Code        string `json:"code,omitempty"` // Código de embed del video
}

// Videos contiene la colección de servidores de video organizados por idioma.
// La estructura refleja el formato JSON del sitio AnimeFlv.
type Videos struct {
	SUB []VideoServer `json:"SUB"` // Servidores de video subtitulados
}
