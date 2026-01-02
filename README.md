# 🎌 Anime API

Librería Go de alto rendimiento para obtener información de animes desde AnimeFlv con caché distribuido opcional.

## 📋 Descripción

**Anime API** te permite buscar animes, obtener información detallada (sinopsis, géneros, estado, episodios, animes relacionados) y conseguir enlaces de episodios desde servicios externos (Mega, Zippyshare, StreamSB, etc.) disponibles en AnimeFlv.

Incluye **caché distribuido opcional** con Valkey/Redis para optimizar consultas recurrentes, reduciendo tiempos de respuesta de 2-3 segundos a <1ms.

### 🌟 Características

- 🔍 **Búsqueda de animes** por nombre con paginación
- 📖 **Información completa** - Sinopsis, géneros, estado, episodios, animes relacionados
- 🎬 **Enlaces de episodios** - URLs de servicios externos (Mega, Zippyshare, StreamSB, etc.)
- 🎥 **Extracción de URLs** - Obtén URLs directas de reproducción (⚡ NEW - StreamTape soportado, más servicios próximamente)
- 📺 **Animes recientes** - Últimos agregados al sitio
- 🆕 **Episodios recientes** - Últimos episodios publicados
- 💾 **Caché opcional** - Configurable, desactivable, TTL personalizable
- 🚀 **Alto rendimiento** - < 1ms en consultas cacheadas

---

## 📦 Instalación

```bash
go get github.com/dst3v3n/api-anime
```

### Prerrequisitos

- **Go 1.25.3+**
- **Valkey/Redis** (opcional, solo si quieres usar caché)

---

## 🚀 Inicio Rápido

### Sin Caché (Más Simple)

```go
package main

import (
    "context"
    "fmt"
    "github.com/dst3v3n/api-anime"
)

func main() {
    // Usa configuración por defecto (caché desactivado)
    service := apianime.NewAnimeFlv()
    ctx := context.Background()
    
    resultados, err := service.SearchAnime(ctx, "One Piece", 1)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Encontrados %d animes\n", len(resultados.Animes))
}
```

### Con Caché (Recomendado)

**1. Iniciar Valkey/Redis:**

```bash
docker run -d -p 6379:6379 valkey/valkey:latest
```

**2. Configurar y usar:**

```go
package main

import (
    "context"
    "github.com/dst3v3n/api-anime"
    "github.com/dst3v3n/api-anime/config"
)

func main() {
    // Activar caché programáticamente
    cfg := config.NewConfigWithDefaults().
        WithEnableCache(true).           // Activar
        WithCacheHost("localhost").      // Host
        WithCachePort(6379).             // Puerto
        WithCacheTTL(60)                 // 60 minutos (1 hora)
    
    config.InitConfig(cfg)
    
    service := apianime.NewAnimeFlv()
    ctx := context.Background()
    
    // Primera búsqueda: ~2s (scraping)
    resultados, _ := service.SearchAnime(ctx, "Naruto", 1)
    
    // Segunda búsqueda: <1ms (desde caché!)
    resultados, _ = service.SearchAnime(ctx, "Naruto", 1)
}
```

---

## 📚 Referencia API

### SearchAnime

Busca animes por nombre con paginación.

```go
SearchAnime(ctx context.Context, anime string, page uint) (AnimeResponse, error)
```

**Ejemplo:**

```go
resultados, err := service.SearchAnime(ctx, "Naruto", 1)
if err != nil {
    log.Fatal(err)
}

for _, anime := range resultados.Animes {
    fmt.Printf("%s - ⭐%.1f\n", anime.Title, anime.Punctuation)
}
```

**Retorna:**

```go
import "github.com/dst3v3n/api-anime/types"

type AnimeResponse struct {
    Animes     []types.AnimeStruct
    TotalPages uint
}

type AnimeStruct struct {
    ID          string        // "naruto-shippuden"
    Title       string        // "Naruto Shippuden"
    Sinopsis    string
    Type        CategoryAnime // Anime, OVA, Pelicula, Especial
    Punctuation float64       // 0-10
    Image       string        // URL
}
```

Disponible en: `types.AnimeResponse` y `types.AnimeStruct`

---

### AnimeInfo

Información completa de un anime.

```go
AnimeInfo(ctx context.Context, idAnime string) (AnimeInfoResponse, error)
```

**Ejemplo:**

```go
info, _ := service.AnimeInfo(ctx, "one-piece-tv")

fmt.Println("Estado:", info.Status)          // "En Emision" / "Finalizado"
fmt.Println("Géneros:", info.Genres)
fmt.Println("Próximo ep:", info.NextEpisode)
fmt.Println("Total eps:", len(info.Episodes))

// Animes relacionados
for _, rel := range info.AnimeRelated {
    fmt.Printf("- %s (%s)\n", rel.Title, rel.Category)
}
```

**Retorna:**

```go
import "github.com/dst3v3n/api-anime/types"

type AnimeInfoResponse struct {
    AnimeStruct                   // Info básica
    AnimeRelated []AnimeRelated   // Secuelas, precuelas
    Genres       []string
    Status       StatusAnime      // "En Emision" / "Finalizado"
    NextEpisode  string
    Episodes     []int            // [1, 2, 3, ..., 1150]
}
```

Disponible en: `types.AnimeInfoResponse` y `types.AnimeRelated`

---

### Links

Obtiene los enlaces de descarga/streaming de un episodio desde diferentes servicios externos (Mega, Zippyshare, StreamSB, etc.).

```go
Links(ctx context.Context, idAnime string, episode uint) (LinkResponse, error)
```

**Ejemplo:**

```go
links, _ := service.Links(ctx, "one-piece-tv", 1150)

// Mostrar todos los servidores disponibles
for _, link := range links.Link {
    fmt.Printf("Servidor: %s\n", link.Server)  // "Mega", "Zippyshare", etc.
    fmt.Printf("URL: %s\n", link.URL)          // Enlace directo al servicio
    fmt.Println("---")
}
```

**Retorna:**

```go
import "github.com/dst3v3n/api-anime/types"

type LinkResponse struct {
    ID      string
    Title   string
    Episode uint
    Link    []types.LinkSource  // Enlaces de servicios externos
}

type LinkSource struct {
    Server string // Nombre del servicio: "Mega", "Zippyshare", "StreamSB", etc.
    URL    string // URL directa al servicio externo
    Code   string // Código embed (si aplica)
}
```

Disponible en: `types.LinkResponse` y `types.LinkSource`

---

### ExtractUrl ⚡ NUEVO

Extrae la URL directa de reproducción desde una página embebida de video. **Actualmente disponible solo para StreamTape**, con soporte para más servicios próximamente.

```go
ExtractUrl(ctx context.Context, url string) (string, error)
```

**Ejemplo:**

```go
import (
    "github.com/dst3v3n/api-anime/extract"
)

// URL del reproductor embebido de StreamTape
embedURL := "https://streamtape.com/e/PWw1erZpe1FG87/"

// Extraer URL directa
videoURL, err := extract.ExtractUrl(ctx, embedURL)
if err != nil {
    log.Fatal(err)
}

fmt.Println("URL del video:", videoURL)
// Output: https://streamtape.com/get_video?id=...
```

**Características:**

- 🎬 Extrae URLs directas desde reproductores embebidos
- 🤖 Usa automatización de navegador (Chromedp/Chrome headless)
- ⏱️ Tiempo de extracción: ~3-5 segundos
- 🔄 Requiere Chrome/Chromium instalado en el sistema

**Servicios Soportados:**

| Servicio | Estado | Notas |
|----------|--------|-------|
| StreamTape | ✅ Disponible | Actualmente soportado |
| Mega | ⏳ Próximamente | En desarrollo |
| Zippyshare | ⏳ Próximamente | En desarrollo |
| Google Drive | ⏳ Próximamente | En desarrollo |

**Requisitos:**

```bash
# Chrome o Chromium debe estar instalado en el sistema
# En Linux
sudo apt-get install chromium-browser

# En macOS
brew install chromium

# En Windows
# Descargar desde: https://www.chromium.org/
```

**Uso Completo - Obtener enlaces y extraer URLs:**

```go
package main

import (
    "context"
    "fmt"
    "github.com/dst3v3n/api-anime"
    "github.com/dst3v3n/api-anime/extract"
)

func main() {
    service := apianime.NewAnimeFlv()
    ctx := context.Background()
    
    // 1. Obtener enlaces del episodio
    links, _ := service.Links(ctx, "one-piece-tv", 1150)
    
    // 2. Buscar servidor StreamTape
    for _, link := range links.Link {
        if link.Server == "StreamTape" {
            fmt.Printf("Encontrado: %s\n", link.URL)
            
            // 3. Extraer URL directa
            videoURL, err := extract.ExtractUrl(ctx, link.URL)
            if err != nil {
                fmt.Println("Error:", err)
                continue
            }
            
            fmt.Println("URL del video:", videoURL)
            break
        }
    }
}
```

---

### RecentAnime

Últimos animes agregados al sitio.

```go
RecentAnime(ctx context.Context) ([]AnimeStruct, error)
```

**Ejemplo:**

```go
recientes, _ := service.RecentAnime(ctx)

for _, anime := range recientes[:5] {
    fmt.Println("-", anime.Title)
}
```

---

### RecentEpisode

Últimos episodios publicados.

```go
RecentEpisode(ctx context.Context) ([]EpisodeListResponse, error)
```

**Ejemplo:**

```go
episodios, _ := service.RecentEpisode(ctx)

for _, ep := range episodios[:5] {
    fmt.Printf("%s - Ep. %d\n", ep.Title, ep.Episode)
}
```

---

## 💡 Casos de Uso

### Buscar y explorar animes

```go
// Buscar por nombre
resultados, _ := service.SearchAnime(ctx, "Attack on Titan", 1)

// Ver detalles
info, _ := service.AnimeInfo(ctx, resultados.Animes[0].ID)

// Explorar relacionados
for _, rel := range info.AnimeRelated {
    fmt.Printf("- %s (%s)\n", rel.Title, rel.Category)
}
```

### Obtener enlaces de todos los episodios

```go
info, _ := service.AnimeInfo(ctx, "shingeki-no-kyojin")

for _, ep := range info.Episodes {
    links, _ := service.Links(ctx, "shingeki-no-kyojin", uint(ep))
    fmt.Printf("Ep.%d tiene %d servicios disponibles:\n", ep, len(links.Link))
    
    // Mostrar cada servicio
    for _, link := range links.Link {
        fmt.Printf("  - %s: %s\n", link.Server, link.URL)
    }
}
```

### Monitorear nuevos episodios

```go
episodios, _ := service.RecentEpisode(ctx)

for _, ep := range episodios {
    fmt.Printf("[NUEVO] %s - Cap. %s\n", ep.Title, ep.Chapter)
}
```

---

## 🔧 Configuración

### Opción 1: Variables de Entorno (.env)

```bash
# .env
CACHE_ENABLED=true
CACHE_HOST=localhost
CACHE_PORT=6379
CACHE_DB=0
CACHE_TTL=60    # minutos
```

```go
// Carga automática
service := apianime.NewAnimeFlv()
```

### Opción 2: Configuración Programática (Recomendado)

```go
import "github.com/dst3v3n/api-anime/config"

// Builder pattern
cfg := config.NewConfigWithDefaults().
    WithEnableCache(true).              // Activar caché
    WithCacheHost("redis.prod.com").    // Host
    WithCachePort(6380).                // Puerto
    WithCachePassword("secret").        // Contraseña
    WithCacheTTL(120)                   // 2 horas (en minutos)

config.InitConfig(cfg)
service := apianime.NewAnimeFlv()
```

### Opción 3: Desde archivo .env personalizado

```go
cfg, err := config.NewConfigFromEnvPath("/custom/.env")
if err != nil {
    panic(err)
}
config.InitConfig(cfg)
```

### Configuración Detallada

| Método | Tipo | Default | Descripción |
|--------|------|---------|-------------|
| `WithEnableCache(bool)` | bool | false | Activar/desactivar caché |
| `WithCacheHost(string)` | string | localhost | Host Valkey/Redis |
| `WithCachePort(int)` | int | 6379 | Puerto (1-65535) |
| `WithCachePassword(string)` | string | "" | Contraseña (opcional) |
| `WithCacheDB(int)` | int | 0 | Base datos (0-15) |
| `WithCacheTTL(int)` | int | 60 | TTL en **minutos** |

### Ejemplos de Configuración

**Desarrollo local sin caché:**

```go
cfg := config.NewConfigWithDefaults()
// No necesitas configurar nada más
```

**Desarrollo con caché local:**

```go
cfg := config.NewConfigWithDefaults().
    WithEnableCache(true)
```

**Producción con Redis:**

```go
cfg := config.NewConfigWithDefaults().
    WithEnableCache(true).
    WithCacheHost("redis-prod.example.com").
    WithCachePort(6380).
    WithCachePassword(os.Getenv("REDIS_PASSWORD")).
    WithCacheTTL(30)  // 30 minutos
```

**Múltiples entornos:**

```go
func newService(env string) *apianime.AnimeFlv {
    var cfg *config.Config
    
    switch env {
    case "production":
        cfg = config.NewConfigWithDefaults().
            WithEnableCache(true).
            WithCacheHost("redis.prod.com").
            WithCacheTTL(60)  // 1 hora
    case "development":
        cfg = config.NewConfigWithDefaults().
            WithEnableCache(false)
    default:
        cfg = config.NewConfigWithDefaults()
    }
    
    config.InitConfig(cfg)
    return apianime.NewAnimeFlv()
}
```

---

## 💾 Sistema de Caché

### ¿Qué se cachea?

| Operación | Clave | TTL Default |
|-----------|-------|-------------|
| SearchAnime | `search-anime-{nombre}-page-{N}` | 15m |
| AnimeInfo | `anime-info-{id}` | 15m |
| Links | `links-{id}-{episodio}` | 15m |
| RecentAnime | `recent-anime` | 15m |
| RecentEpisode | `recent-episode` | 15m |

### Performance

| Operación | Sin Caché | Con Caché | Mejora |
|-----------|-----------|-----------|--------|
| SearchAnime | 2.5s | 0.8ms | **3100x** |
| AnimeInfo | 1.8s | 0.6ms | **3000x** |
| Links | 1.5s | 0.5ms | **3000x** |

### Activar/Desactivar en Tiempo Real

```go
// Desactivar caché temporalmente
cfg := config.NewConfigWithDefaults().WithEnableCache(false)
config.InitConfig(cfg)

// Búsqueda sin caché
service.SearchAnime(ctx, "Naruto", 1)

// Reactivar caché
cfg.WithEnableCache(true)
config.InitConfig(cfg)
```

---

## ❓ FAQ

**¿Necesito Valkey/Redis obligatoriamente?**  
No, el caché está desactivado por defecto. Funciona perfectamente sin él.

**¿Cómo activo el caché?**  

```go
cfg := config.NewConfigWithDefaults().WithEnableCache(true)
config.InitConfig(cfg)
```

**¿Puedo cambiar el TTL?**  
Sí, usa `WithCacheTTL(minutos)`:

```go
cfg.WithCacheTTL(120)  // 2 horas
```

**¿Funciona con Redis en lugar de Valkey?**  
Sí, son 100% compatibles. Usa los mismos métodos de configuración.

**¿Los enlaces caducan?**  
Sí, algunos servidores tienen enlaces temporales. Por eso el caché tiene TTL de 15 minutos por defecto.

**¿Puedo usar en producción?**  
Sí, pero el scraping depende de la estructura del sitio. Monitorea cambios regularmente.

---

## 🧪 Testing

```bash
# Tests completos
go test ./...

# Con cobertura
go test ./... -cover

# Tests específicos
go test ./internal/adapters/scrapers/animeflv -v
```

---

## ⚠️ Aviso Legal

**Solo para fines educativos**. Respeta los términos de servicio de AnimeFlv.

**Obligaciones:**

- ✅ Respeta `robots.txt`
- ✅ Usa para proyectos personales/educativos
- ✅ Cita la fuente (AnimeFlv)
- ✅ Implementa rate limiting

**Prohibido:**

- ❌ Comercialización sin permiso
- ❌ Ataques DDoS o sobrecarga
- ❌ Distribución sin atribución

---

## 📄 Licencia

MIT - Ver [LICENSE](LICENSE) para detalles.

---

## 👤 Autor

**Steven** ([@dst3v3n](https://github.com/dst3v3n))

---

## 🤝 Contribuir

¡Contribuciones bienvenidas!

- 🐛 Bugs: [Issues](../../issues)
- 💡 Features: [Discussions](../../discussions)
- 🔧 Código: [Pull Request](../../pulls)

---

## 📞 Soporte

- **GitHub**: [@dst3v3n](https://github.com/dst3v3n)
- **Issues**: [GitHub Issues](../../issues)

---

**Made with ❤️ by Steven**
