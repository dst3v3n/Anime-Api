// Package animeflv - mapper.go
// Este archivo implementa el componente Mapper que se encarga de transformar
// datos crudos extraídos del HTML en estructuras DTO (Data Transfer Objects) bien definidas.
// Convierte los datos primitivos en objetos de dominio utilizables por la aplicación,
// proporcionando una capa de abstracción entre el scraping y la lógica de negocio.
// Los mappers normalizan y estructuran los datos para garantizar consistencia
// y facilitar la validación en las capas superiores.
package animeflv

import "github.com/dst3v3n/api-anime/internal/domain/dto"

// Maper es el componente encargado de transformar datos primitivos extraídos del HTML
// en objetos DTO (Data Transfer Objects) bien definidos. Implementa el patrón Mapper
// para desacoplar la lógica de scraping de la estructura de datos de dominio.
type Maper struct{}

// NewMaper crea una nueva instancia del mapper.
// Retorna un pointer a Maper listo para ser usado en transformaciones de datos.
func NewMaper() *Maper {
	return &Maper{}
}

// ToAnime transforma datos básicos de anime en un DTO AnimeStruct.
// Normaliza tipos primitivos en una estructura bien definida con tipos específicos del dominio.
// Parámetros:
//   - ID: identificador único del anime (ej: "one-piece-tv")
//   - Title: título del anime
//   - Sinopsis: descripción/sinopsis del anime
//   - Tipo: categoría como string (ej: "Anime", "Película", "OVA")
//   - Punctuation: calificación numérica del anime
//   - Image: URL de la imagen/carátula del anime
//
// Retorna: estructura AnimeStruct con tipos enumerados del dominio
func (m *Maper) ToAnime(ID string, Title string, Sinopsis string, Tipo string, Punctuation float64, Image string) dto.AnimeStruct {
	return dto.AnimeStruct{
		ID:          ID,
		Title:       Title,
		Sinopsis:    Sinopsis,
		Type:        dto.CategoryAnime(Tipo),
		Punctuation: Punctuation,
		Image:       Image,
	}
}

// ToAnimeInfo transforma datos completos de anime en un DTO AnimeInfoResponse.
// Combina información básica (heredada de AnimeStruct) con datos adicionales como
// géneros, episodios disponibles, animes relacionados y estado de emisión.
// Parámetros:
//   - ID: identificador único del anime
//   - Title: título del anime
//   - Sinopsis: descripción del anime
//   - Tipo: categoría como string
//   - Punctuation: calificación numérica
//   - Image: URL de la imagen
//   - AnimeRelated: slice de animes relacionados (secuelas, precuelas, spin-offs)
//   - Generos: slice de géneros del anime
//   - Estado: estado de emisión como string ("En Emision" o "Finalizado")
//   - Episodes: slice de números de episodios disponibles
//   - NextEpisode: fecha del próximo episodio a emitirse
//
// Retorna: estructura AnimeInfoResponse con información completa del anime
func (m *Maper) ToAnimeInfo(ID string, Title string, Sinopsis string, Tipo string, Punctuation float64, Image string, AnimeRelated []dto.AnimeRelated, Generos []string, Estado string, Episodes []int, NextEpisode string) dto.AnimeInfoResponse {
	return dto.AnimeInfoResponse{
		AnimeStruct: dto.AnimeStruct{
			ID:          ID,
			Title:       Title,
			Sinopsis:    Sinopsis,
			Type:        dto.CategoryAnime(Tipo),
			Punctuation: Punctuation,
			Image:       Image,
		},
		AnimeRelated: AnimeRelated,
		Genres:       Generos,
		Status:       dto.StatusAnime(Estado),
		NextEpisode:  NextEpisode,
		Episodes:     Episodes,
	}
}

// ToLinks transforma datos de un servidor de video individual en un DTO LinkSource.
// Encapsula la información de un único servidor de reproducción.
// Parámetros:
//   - Server: nombre del servidor de video (ej: "Zippyshare", "Mega", "Google Drive", "StreamTape")
//   - URL: URL del enlace o página del servidor
//   - Code: código de embed o identificador del video en el servidor
//
// Retorna: estructura LinkSource con información de un servidor de reproducción
func (m *Maper) ToLinks(Server string, URL string, Code string) dto.LinkSource {
	return dto.LinkSource{
		Server: Server,
		URL:    URL,
		Code:   Code,
	}
}

// ToLinkEpisode transforma datos de enlaces de episodio en un DTO LinkResponse.
// Agrupa todos los enlaces de diferentes servidores para un episodio específico,
// creando una respuesta completa con información del episodio y todas sus opciones de reproducción.
// Parámetros:
//   - ID: identificador único del anime
//   - Title: título del anime
//   - Episode: número del episodio
//   - Links: slice de LinkSource con opciones de reproducción
//
// Retorna: estructura LinkResponse con información del episodio y sus enlaces
func (m *Maper) ToLinkEpisode(ID string, Title string, Episode uint, Links []dto.LinkSource) dto.LinkResponse {
	return dto.LinkResponse{
		ID:      ID,
		Episode: Episode,
		Title:   Title,
		Link:    Links,
	}
}

// ToRecentEpisode transforma datos de episodio reciente en un DTO EpisodeListResponse.
// Crea una representación resumida de un episodio para mostrar en listados de episodios recientes.
// Incluye información esencial: identificadores, título, número de episodio y imagen.
// Parámetros:
//   - ID: identificador único del anime
//   - Title: título del anime
//   - Chapter: designación del capítulo (ej: "Cap. 1050")
//   - Episode: número del episodio
//   - Image: URL de la imagen/carátula del episodio
//
// Retorna: estructura EpisodeListResponse para uso en listados
func (m *Maper) ToRecentEpisode(ID string, Title string, Chapter string, Episode int, Image string) dto.EpisodeListResponse {
	return dto.EpisodeListResponse{
		ID:      ID,
		Title:   Title,
		Chapter: Chapter,
		Episode: Episode,
		Image:   Image,
	}
}
