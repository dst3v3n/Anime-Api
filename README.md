# 🎌 Anime API

API de scraping para obtener información de animes desde AnimeFlv.

## 📋 Descripción

Esta librería te permite buscar animes, obtener información detallada (sinopsis, géneros, episodios) y conseguir los enlaces de descarga/streaming de cada episodio desde AnimeFlv.

## ✨ Características

- 🔍 **Búsqueda de animes** por nombre con paginación
- 📖 **Información detallada** de cada anime (sinopsis, géneros, estado, episodios disponibles)
- 🎬 **Enlaces de episodios** para ver o descargar de múltiples servidores
- 📺 **Animes recientes** obtén los últimos animes agregados al sitio
- 🆕 **Episodios recientes** listado de los últimos episodios publicados
- 💾 **Caché distribuido** Valkey integrado para optimizar consultas recurrentes
- 🚀 **Alto rendimiento** con scraping optimizado usando goquery y caché
- 🏗️ **Arquitectura hexagonal** diseño limpio y extensible

## 📦 Instalación

### Prerrequisitos

- Go 1.23 o superior
- Valkey/Redis en ejecución (para caché distribuido) - puerto 6379 por defecto

### Instalar la librería

```bash
go get github.com/dst3v3n/api-anime
```

### Dependencias

El proyecto utiliza las siguientes dependencias:
- `github.com/PuerkitoBio/goquery` - Parser HTML para scraping eficiente
- `github.com/valkey-io/valkey-go` - Cliente Valkey para caché distribuido

## 🚀 Uso

### Configuración Inicial

Asegúrate de tener Valkey ejecutándose:

```bash
# Usando Docker (recomendado)
docker run -d -p 6379:6379 valkey/valkey:latest

# O instala Valkey localmente
# Ubuntu/Debian: sudo apt-get install valkey
# macOS: brew install valkey
```

### Inicio rápido

```go
package main

import (
    "fmt"
    "github.com/dst3v3n/api-anime/internal/domain/services/animeflv"
)

func main() {
    // Crear servicio con caché integrado
    service := animeflv.NewAnimeflvService()

    // 1. Buscar anime por nombre
    anime := "One Piece"
    page := uint(1)
    resultados, err := service.SearchAnime(&anime, &page)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    for _, animeResult := range resultados.Animes {
        fmt.Println("Título:", animeResult.Title)
        fmt.Println("ID:", animeResult.ID)
        fmt.Println("Tipo:", animeResult.Type)
        fmt.Println("Puntuación:", animeResult.Punctuation)
        fmt.Println("------------------------------------")
    }

    // 2. Obtener información detallada (con caché)
    idAnime := "one-piece-tv"
    info, err := service.AnimeInfo(&idAnime)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Título:", info.Title)
    fmt.Println("Estado:", info.Status)
    fmt.Println("Géneros:", info.Genres)
    fmt.Println("Episodios disponibles:", len(info.Episodes))
    
    // 3. Obtener enlaces de un episodio específico (con caché)
    episode := uint(1150)
    links, err := service.Links(&idAnime, &episode)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Episodio:", links.Episode)
    for _, link := range links.Link {
        fmt.Println("Servidor:", link.Server)
        fmt.Println("URL:", link.URL)
    }

    // 4. Obtener animes recientes (con caché)
    recientes, err := service.RecentAnime()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Animes recientes:", len(recientes))
    for _, animeRecent := range recientes {
        fmt.Println("-", animeRecent.Title)
    }

    // 5. Obtener episodios recientes (con caché)
    episodios, err := service.RecentEpisode()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Episodios recientes:", len(episodios))
    for _, ep := range episodios {
        fmt.Printf("- %s - Episodio %d\n", ep.Title, ep.Episode)
    }
}
```

## 📚 Referencia de Métodos

### 🔍 SearchAnime

Busca animes por nombre con paginación. **Incluye caché automático de 15 minutos**.

```go
SearchAnime(anime *string, page *uint) (dto.AnimeResponse, error)
```

**Parámetros:**
- `anime` (*string): Nombre del anime a buscar
- `page` (*uint): Número de página (ejemplo: 1, 2, 3, etc.)

**Retorna:**
```go
type AnimeResponse struct {
    ID          string        // ID único del anime (ej: "one-piece-tv")
    Title       string        // Título del anime
    Sinopsis    string        // Breve descripción
    Type        CategoryAnime // Tipo: Anime, OVA, Película, Especial
    Punctuation float64       // Puntuación (0.0 - 5.0)
    Image       string        // URL de la imagen de portada
}
```

**Ejemplo:**
```go
anime := "Naruto"
page := uint(1)
resultados, err := service.SearchAnime(&anime, &page)
```

> **Nota:** Los resultados se cachean automáticamente en Valkey por 15 minutos. Búsquedas posteriores de la misma query retornarán datos del caché instantáneamente.

---

### 📖 AnimeInfo

Obtiene información detallada de un anime específico. **Incluye caché automático de 15 minutos**.

```go
AnimeInfo(idAnime *string) (dto.AnimeInfoResponse, error)
```

**Parámetros:**
- `idAnime` (*string): ID del anime (obtenido de SearchAnime)

**Retorna:**
```go
type AnimeInfoResponse struct {
    AnimeResponse                        // Hereda los campos básicos de AnimeResponse
    AnimeRelated []AnimeRelated          // Animes relacionados (precuelas, secuelas, etc.)
    Genres       []string                // Lista de géneros
    Status       StatusAnime             // "En Emision" o "Finalizado"
    NextEpisode  string                  // Fecha del próximo episodio (si está en emisión)
    Episodes     []int                   // Lista de números de episodios disponibles
}

type AnimeRelated struct {
    ID       string  // ID del anime relacionado
    Title    string  // Título del anime relacionado
    Category string  // Tipo de relación (Secuela, Precuela, etc.)
}
```

**Ejemplo:**
```go
id := "naruto-shippuden"
info, err := service.AnimeInfo(&id)
```

> **Nota:** Los detalles se cachean automáticamente por 15 minutos bajo la clave `anime-info-{id}`.

---

### 🔗 Links

Obtiene los enlaces de descarga/streaming de un episodio específico. **Incluye caché automático de 15 minutos**.

```go
Links(idAnime *string, episode *uint) (dto.LinkResponse, error)
```

**Parámetros:**
- `idAnime` (*string): ID del anime
- `episode` (*uint): Número del episodio

**Retorna:**
```go
type LinkResponse struct {
    ID      string       // ID del anime
    Title   string       // Título del episodio
    Episode int          // Número del episodio
    Link    []LinkSource // Enlaces de diferentes servidores
}

type LinkSource struct {
    Server string // Nombre del servidor (ej: "Zippyshare", "Mega", etc.)
    URL    string // URL de descarga/streaming
    Code   string // Código de embed del video
}
```

**Ejemplo:**
```go
id := "naruto-shippuden"
episode := uint(1)
links, err := service.Links(&id, &episode)
```

> **Nota:** Los enlaces se cachean automáticamente por 15 minutos bajo la clave `links-{id}-{episode}`.

---

### 📺 RecentAnime

Obtiene la lista de animes recientemente agregados al sitio. **Incluye caché automático de 15 minutos**.

```go
RecentAnime() ([]dto.AnimeStruct, error)
```

**Retorna:**
Lista de `AnimeStruct` con los animes recientes.

**Ejemplo:**
```go
recientes, err := service.RecentAnime()
for _, anime := range recientes {
    fmt.Printf("%s - %s\n", anime.Title, anime.Type)
}
```

> **Nota:** Los animes recientes se cachean bajo la clave `recent-anime` por 15 minutos.

---

### 🆕 RecentEpisode

Obtiene la lista de episodios recientemente publicados. **Incluye caché automático de 15 minutos**.

```go
RecentEpisode() ([]dto.EpisodeListResponse, error)
```

**Retorna:**
```go
type EpisodeListResponse struct {
    ID      string // ID del anime
    Title   string // Título del anime
    Chapter string // Texto descriptivo del capítulo
    Episode int    // Número del episodio
    Image   string // URL de la imagen del episodio
}
```

**Ejemplo:**
```go
episodios, err := service.RecentEpisode()
for _, ep := range episodios {
    fmt.Printf("%s - Episodio %d\n", ep.Title, ep.Episode)
}
```

> **Nota:** Los episodios recientes se cachean bajo la clave `recent-episode` por 15 minutos.

## 💡 Casos de Uso

### Buscar y listar animes
```go
// Buscar "Attack on Titan" en la primera página (con caché)
title := "Attack on Titan"
page := uint(1)
resultados, _ := service.SearchAnime(&title, &page)
for _, anime := range resultados.Animes {
    fmt.Printf("%s (%s) - ⭐%.1f\n", anime.Title, anime.Type, anime.Punctuation)
}
// La segunda búsqueda de "Attack on Titan" página 1 será instantánea (desde caché)
```

### Obtener todos los episodios de un anime
```go
id := "shingeki-no-kyojin"
info, _ := service.AnimeInfo(&id)
fmt.Printf("Estado: %s\n", info.Status)
fmt.Printf("Total de episodios: %d\n", len(info.Episodes))

// Obtener enlaces de todos los episodios (con caché)
for _, ep := range info.Episodes {
    episode := uint(ep)
    links, _ := service.Links(&id, &episode)
    fmt.Printf("Episodio %d tiene %d enlaces\n", ep, len(links.Link))
}
```

### Verificar nuevos episodios
```go
id := "one-piece-tv"
info, _ := service.AnimeInfo(&id)
if info.Status == "En Emision" {
    fmt.Println("Próximo episodio:", info.NextEpisode)
    ultimoEp := info.Episodes[len(info.Episodes)-1]
    fmt.Println("Último episodio disponible:", ultimoEp)
}
```

### Monitorear animes y episodios recientes
```go
// Ver qué animes nuevos se agregaron (con caché)
recientes, _ := service.RecentAnime()
fmt.Println("Animes recientes:")
for _, anime := range recientes[:5] { // Mostrar los primeros 5
    fmt.Printf("- %s (⭐%.1f)\n", anime.Title, anime.Punctuation)
}

// Ver qué episodios nuevos salieron hoy (con caché)
episodios, _ := service.RecentEpisode()
fmt.Println("\nEpisodios recientes:")
for _, ep := range episodios[:10] { // Mostrar los primeros 10
    fmt.Printf("- %s - Ep. %d\n", ep.Title, ep.Episode)
}
// Las siguientes llamadas retornarán datos del caché en < 1ms
```

### Explorar animes relacionados
```go
id := "naruto"
info, _ := service.AnimeInfo(&id)
fmt.Println("Animes relacionados con Naruto:")
for _, related := range info.AnimeRelated {
    fmt.Printf("- %s (%s)\n", related.Title, related.Category)
}
```

## 🔧 Manejo de Errores

Todos los métodos retornan un error que debes manejar:

```go
anime := "Naruto"
page := uint(1)
resultados, err := service.SearchAnime(&anime, &page)
if err != nil {
    log.Fatal("Error en la búsqueda:", err)
}
```

**Errores comunes:**
- Anime no encontrado
- Episodio no disponible
- Problemas de conexión con el sitio web
- Cambios en la estructura del sitio (requiere actualización de la librería)

## 💾 Sistema de Caché

La API incluye **caché distribuido integrado** usando Valkey (alternativa a Redis):

### Características del Caché

- **Automático**: Se aplica automáticamente a todas las operaciones sin configuración adicional
- **TTL configurable**: 15 minutos por defecto para búsquedas y contenido reciente
- **Distribuido**: Perfecto para aplicaciones con múltiples instancias
- **Transparente**: Los desarrolladores no necesitan gestionar el caché manualmente
- **Optimizado**: Las búsquedas repetidas retornan resultados en < 1ms desde caché

### Claves de Caché

| Operación | Clave Caché | TTL |
|-----------|------------|-----|
| `SearchAnime("one-piece", 1)` | `search-anime-one-piece-page-1` | 15 min |
| `AnimeInfo("one-piece-tv")` | `anime-info-one-piece-tv` | 15 min |
| `Links("one-piece-tv", 1150)` | `links-one-piece-tv-1150` | 15 min |
| `RecentAnime()` | `recent-anime` | 15 min |
| `RecentEpisode()` | `recent-episode` | 15 min |

### Ventajas del Caché

✅ **Rendimiento**: Búsquedas repetidas son instantáneas
✅ **Escalabilidad**: Soporta múltiples instancias de la aplicación
✅ **Confiabilidad**: Reduce carga en el servidor remoto
✅ **Experiencia**: Aplicación más responsiva

## 📊 Tipos de Datos

### CategoryAnime
Tipos de contenido disponibles:
- `Anime` - Series de anime regulares
- `Ova` - Original Video Animation
- `Pelicula` - Películas de anime
- `Especial` - Episodios especiales

### StatusAnime
Estado de emisión:
- `En Emision` - Anime actualmente en emisión
- `Finalizado` - Anime completado

## 🏗️ Arquitectura

El proyecto sigue los principios de **Arquitectura Hexagonal (Ports & Adapters)**:

```
├── internal/
│   ├── ports/              # Interfaces (contratos)
│   │   ├── scraper.go      # ScraperPort - Define operaciones de scraping
│   │   └── mapper.go       # MapperPort - Define operaciones de mapeo
│   │
│   └── adapters/           # Implementaciones
│       └── scrapers/
│           ├── dto/        # Data Transfer Objects
│           │   ├── anime.go
│           │   ├── anime_info.go
│           │   ├── episodeList.go
│           │   └── link.go
│           │
│           └── animeflv/   # Implementación para AnimeFlv
│               ├── client.go       # Cliente HTTP y orquestación
│               ├── html_parser.go  # Parsing de HTML con goquery
│               ├── script_parser.go# Extracción de datos de scripts JS
│               ├── mapper.go       # Transformación de datos a DTOs
│               ├── helper.go       # Funciones auxiliares
│               └── models.go       # Modelos internos
```

### Ventajas de esta arquitectura:
- ✅ **Desacoplamiento**: La lógica de negocio no depende de detalles de implementación
- ✅ **Testeable**: Fácil crear mocks de las interfaces para testing
- ✅ **Extensible**: Agregar nuevos scrapers sin modificar código existente
- ✅ **Mantenible**: Cada componente tiene una responsabilidad clara

## 🧪 Testing

El proyecto incluye tests unitarios e integración:

```bash
# Ejecutar tests unitarios
go test ./test/unit/...

# Ejecutar tests de integración
go test ./test/integration/...

# Ejecutar todos los tests
go test ./...
```

## ❓ FAQ

**¿Puedo usar esto en producción?**
Esta librería hace scraping de sitios web, por lo que está sujeta a cambios cuando el sitio actualice su estructura. Úsala bajo tu propio riesgo.

**¿Qué tan rápido es?**
La velocidad depende de la conexión y la carga del sitio web. Cada petición toma entre 1-3 segundos aproximadamente. Usa goquery que es más eficiente que colly para este caso de uso.

**¿Funciona con otros sitios además de AnimeFlv?**
Actualmente solo soporta AnimeFlv. Sin embargo, gracias a la arquitectura hexagonal, es fácil agregar nuevos scrapers implementando la interfaz `ScraperPort`.

**¿Los enlaces de descarga caducan?**
Los enlaces son obtenidos en tiempo real del sitio, pero algunos servidores pueden tener enlaces temporales que expiran.

**¿Por qué goquery en lugar de colly?**
Para este proyecto, goquery proporciona todo lo necesario con mejor control sobre las peticiones HTTP y un parsing HTML más directo. Es más ligero y suficiente para las necesidades del scraper.

## 🛠️ Tecnologías Utilizadas

- **Go 1.23+** - Lenguaje de programación
- **goquery** - Parsing y manipulación de HTML
- **Valkey** - Caché distribuido de alto rendimiento
- **net/http** - Cliente HTTP estándar de Go
- **encoding/json** - Procesamiento de datos JSON embebidos

## ⚠️ Aviso Legal

Este proyecto es **solo para fines educativos**. El scraping de sitios web debe hacerse respetando los términos de servicio del sitio. El autor no se hace responsable del uso indebido de esta herramienta.

**Recomendaciones:**
- Implementa rate limiting para no sobrecargar el servidor
- Respeta el archivo `robots.txt` del sitio
- No uses esto para actividades comerciales sin permiso
- Considera usar APIs oficiales cuando estén disponibles

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 👤 Autor

**Steven** - [@dst3v3n](https://github.com/dst3v3n)

## 🤝 Contribuciones

Las contribuciones son bienvenidas! Si quieres agregar soporte para otro sitio de anime o mejorar el código:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agrega nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

## 📝 Roadmap

- [x] Implementar caché distribuido con Valkey
- [ ] Agregar soporte para más sitios de anime
- [ ] Agregar rate limiting configurable
- [ ] Crear CLI para uso desde terminal
- [ ] Agregar más tests de cobertura
- [ ] Documentación con ejemplos adicionales
- [ ] Configurar TTL del caché por operación
- [ ] Endpoint para limpiar caché manualmente

---

⭐ **¿Te resultó útil?** Dale una estrella al repositorio
