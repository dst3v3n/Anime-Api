# 🎌 Anime API

API de scraping para obtener información de animes desde AnimeFlv.

## 📋 Descripción

Esta librería te permite buscar animes, obtener información detallada (sinopsis, géneros, episodios) y conseguir los enlaces de descarga/streaming de cada episodio desde AnimeFlv.

## ✨ Características

- 🔍 **Búsqueda de animes** por nombre con paginación
- 📖 **Información detallada** de cada anime (sinopsis, géneros, estado, episodios disponibles)
- � **Enlaces de episodios** para ver o descargar de múltiples servidores
- 🚀 **Alto rendimiento** con scraping optimizado

## 📦 Instalación

### Prerrequisitos

- Go 1.23 o superior

### Instalar la librería

```bash
go get github.com/dst3v3n/api-anime
```

## 🚀 Uso

### Inicio rápido

```go
package main

import (
    "fmt"
    "github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
)

func main() {
    // Crear cliente
    scraper := animeflv.NewClient()

    // 1. Buscar anime por nombre
    resultados, err := scraper.SearchAnime("One Piece", "1")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    for _, anime := range resultados {
        fmt.Println("Título:", anime.Title)
        fmt.Println("ID:", anime.ID)
        fmt.Println("Tipo:", anime.Tipo)
        fmt.Println("Puntuación:", anime.Puctuation)
        fmt.Println("------------------------------------")
    }

    // 2. Obtener información detallada
    info, err := scraper.AnimeInfo("one-piece-tv")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Título:", info.Title)
    fmt.Println("Estado:", info.Estado)
    fmt.Println("Géneros:", info.Generos)
    fmt.Println("Episodios disponibles:", info.Episodes)
    
    // 3. Obtener enlaces de un episodio específico
    links, err := scraper.GetLinks("one-piece-tv", 1145)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Episodio:", links.Episode)
    fmt.Println("Enlaces disponibles:", links.Links)
}
```

## 📚 Referencia de Métodos

### 🔍 SearchAnime

Busca animes por nombre con paginación.

```go
SearchAnime(anime string, page string) ([]dto.SearchAnimeResponse, error)
```

**Parámetros:**
- `anime` (string): Nombre del anime a buscar
- `page` (string): Número de página (ejemplo: "1", "2", "3", etc.)

**Retorna:**
```go
type SearchAnimeResponse struct {
    ID         string        // ID único del anime (ej: "one-piece-tv")
    Title      string        // Título del anime
    Sipnopsis  string        // Breve descripción
    Tipo       CategoryAnime // Tipo: Anime, OVA, Película
    Puctuation float64       // Puntuación (0.0 - 5.0)
    Image      string        // URL de la imagen de portada
}
```

**Ejemplo:**
```go
resultados, err := scraper.SearchAnime("Naruto", "1")
```

---

### 📖 AnimeInfo

Obtiene información detallada de un anime específico.

```go
AnimeInfo(idAnime string) (dto.ResponseAnimeInfo, error)
```

**Parámetros:**
- `idAnime` (string): ID del anime (obtenido de SearchAnime)

**Retorna:**
```go
type ResponseAnimeInfo struct {
    ID           string
    Title        string
    Sipnopsis    string
    Tipo         CategoryAnime
    Puctuation   float64
    Image        string
    AnimeRelated []AnimeRelated // Animes relacionados (precuelas, secuelas, etc.)
    Generos      []string       // Lista de géneros
    Estado       StatusAnime    // "En Emision" o "Finalizado"
    NextEpisode  string         // Fecha del próximo episodio (si está en emisión)
    Episodes     []int          // Lista de números de episodios disponibles
}
```

**Ejemplo:**
```go
info, err := scraper.AnimeInfo("naruto-shippuden")
```

---

### 🔗 GetLinks

Obtiene los enlaces de descarga/streaming de un episodio específico.

```go
GetLinks(idAnime string, episode int) (dto.ResponseGetLinks, error)
```

**Parámetros:**
- `idAnime` (string): ID del anime
- `episode` (int): Número del episodio

**Retorna:**
```go
type ResponseGetLinks struct {
    ID      string
    Title   string
    Episode int
    Links   []Link // Enlaces de diferentes servidores
}

type Link struct {
    Server string // Nombre del servidor (ej: "Zippyshare", "Mega", etc.)
    URL    string // URL de descarga/streaming
}
```

**Ejemplo:**
```go
links, err := scraper.GetLinks("naruto-shippuden", 1)
```

## 💡 Casos de Uso

### Buscar y listar animes
```go
// Buscar "Attack on Titan" en la primera página
resultados, _ := scraper.SearchAnime("Attack on Titan", "1")
for _, anime := range resultados {
    fmt.Printf("%s (%s) - ⭐%.1f\n", anime.Title, anime.Tipo, anime.Puctuation)
}
```

### Obtener todos los episodios de un anime
```go
info, _ := scraper.AnimeInfo("shingeki-no-kyojin")
fmt.Printf("Estado: %s\n", info.Estado)
fmt.Printf("Total de episodios: %d\n", len(info.Episodes))

// Descargar todos los episodios
for _, ep := range info.Episodes {
    links, _ := scraper.GetLinks("shingeki-no-kyojin", ep)
    fmt.Printf("Episodio %d tiene %d enlaces\n", ep, len(links.Links))
}
```

### Verificar nuevos episodios
```go
info, _ := scraper.AnimeInfo("one-piece-tv")
if info.Estado == "En Emision" {
    fmt.Println("Próximo episodio:", info.NextEpisode)
    ultimoEp := info.Episodes[len(info.Episodes)-1]
    fmt.Println("Último episodio disponible:", ultimoEp)
}
```

## 🔧 Manejo de Errores

Todos los métodos retornan un error que debes manejar:

```go
resultados, err := scraper.SearchAnime("Naruto", "1")
if err != nil {
    log.Fatal("Error en la búsqueda:", err)
}
```

**Errores comunes:**
- Anime no encontrado
- Episodio no disponible
- Problemas de conexión con el sitio web
- Cambios en la estructura del sitio (requiere actualización de la librería)

## 📊 Tipos de Datos

### CategoryAnime
Tipos de contenido disponibles:
- `Anime` - Series de anime regulares
- `OVA` - Original Video Animation
- `Pelicula` - Películas de anime
- `Especial` - Episodios especiales

### StatusAnime
Estado de emisión:
- `En Emision` - Anime actualmente en emisión
- `Finalizado` - Anime completado

## ❓ FAQ

**¿Puedo usar esto en producción?**
Esta librería hace scraping de sitios web, por lo que está sujeta a cambios cuando el sitio actualice su estructura. Úsala bajo tu propio riesgo.

**¿Qué tan rápido es?**
La velocidad depende de la conexión y la carga del sitio web. Cada petición toma entre 1-3 segundos aproximadamente.

**¿Funciona con otros sitios además de AnimeFlv?**
Actualmente solo soporta AnimeFlv. Otros sitios podrían agregarse en el futuro.

**¿Los enlaces de descarga caducan?**
Los enlaces son obtenidos en tiempo real del sitio, pero algunos servidores pueden tener enlaces temporales que expiran.

## ⚠️ Aviso Legal

Este proyecto es **solo para fines educativos**. El scraping de sitios web debe hacerse respetando los términos de servicio del sitio. El autor no se hace responsable del uso indebido de esta herramienta.

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 👤 Autor

**Steven** - [@dst3v3n](https://github.com/dst3v3n)

---

⭐ **¿Te resultó útil?** Dale una estrella al repositorio
