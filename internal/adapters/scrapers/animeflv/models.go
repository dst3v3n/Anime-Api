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
// Centraliza todos los endpoints necesarios para acceder a diferentes secciones del sitio,
// facilitando cambios en la estructura de URLs sin afectar el resto del código.
type Config struct {
	BaseURL       string // URL base del sitio (ej: "https://www3.animeflv.net"). Usada para obtener contenido reciente
	SearchURL     string // URL del endpoint de búsqueda (ej: "https://www3.animeflv.net/browse"). Para búsquedas y listados
	AnimeInfoURL  string // URL del endpoint de información de anime (ej: "https://www3.animeflv.net/anime"). Para detalles completos
	VerEpisodeURL string // URL base para ver episodios (ej: "https://www3.animeflv.net/ver"). Para obtener enlaces de reproducción
}

// ParseResult almacena temporalmente los datos extraídos durante el parsing de información de anime.
// Se utiliza como estructura intermedia (DTO interno) antes de convertir a AnimeInfoResponse.
// Este patrón permite acumular datos de múltiples fuentes (HTML y scripts) antes de la transformación final.
// Campos privados (minúsculas) para encapsulación; solo accesibles dentro del paquete.
type ParseResult struct {
	title        string             // Título del anime extraído del HTML
	category     string             // Categoría/Tipo del anime extraído del HTML (ej: "Anime", "Película", "OVA")
	sipnopsis    string             // Sinopsis o descripción del anime extraída del HTML
	status       string             // Estado de emisión del anime extraído del HTML (ej: "En Emision", "Finalizado")
	image        string             // URL de la imagen/carátula del anime extraída del HTML
	punctuacion  float64            // Calificación del anime extraída del HTML (escala 0-10)
	animeRelated []dto.AnimeRelated // Animes relacionados extraídos del HTML (secuelas, precuelas, spin-offs)
	genres       []string           // Géneros del anime extraídos del HTML
	episodes     []int              // Lista de números de episodios disponibles extraída de scripts JavaScript
	nextEpisode  string             // Fecha del próximo episodio extraída de scripts JavaScript (puede estar vacía si finalizado)
}

// ParseEpisodeLinksResult almacena temporalmente los enlaces extraídos de un episodio.
// Se utiliza como estructura intermedia (DTO interno) antes de convertir a LinkResponse.
// Agrupa información del episodio con todos sus enlaces de reproducción disponibles.
type ParseEpisodeLinksResult struct {
	ID      string           // Identificador único del anime
	Title   string           // Título del anime
	Episode uint             // Número del episodio
	links   []dto.LinkSource // Enlaces de reproducción disponibles (diferentes servidores)
}

// VideoServer representa un servidor de video individual con sus propiedades.
// Se utiliza para deserializar el JSON embebido en los scripts de AnimeFlv.
// Cada servidor de video tiene información sobre su capacidad, tipo y accesibilidad.
type VideoServer struct {
	Server      string `json:"server"`         // Nombre del servidor de video (ej: "Zippyshare", "Mega", "Google Drive", "StreamTape")
	Title       string `json:"title"`          // Título o descripción del servidor (ej: "Zippyshare (HD)")
	Ads         int    `json:"ads"`            // Número aproximado de anuncios en la página del servidor
	URL         string `json:"url,omitempty"`  // URL del enlace directo o página de descarga
	AllowMobile bool   `json:"allow_mobile"`   // Booleano indicando si el servidor es accesible desde dispositivos móviles
	Code        string `json:"code,omitempty"` // Código de embed o identificador del video en el servidor para reproducción integrada
}

// Videos contiene la colección de servidores de video organizados por idioma/versión.
// Refleja exactamente la estructura JSON embebida en los scripts de AnimeFlv.
// Permite deserializar el JSON embebido en las etiquetas <script> de la página de episodios.
type Videos struct {
	SUB []VideoServer `json:"SUB"` // Array de servidores de video subtitulados (SUB = Subtitled/Doblado)
	// Nota: AnimeFlv también puede contener LAT (Latinoamericano) pero actualmente enfocamos en SUB
}
