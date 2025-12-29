# 🤝 Guía de Contribución

¡Gracias por tu interés en contribuir a **Anime API**! Este documento te guiará en el proceso de contribución.

## 📋 Tabla de Contenidos

- [Código de Conducta](#-código-de-conducta)
- [¿Cómo puedo contribuir?](#-cómo-puedo-contribuir)
  - [Reportar Bugs](#-reportar-bugs)
  - [Sugerir Features](#-sugerir-features)
  - [Contribuir Código](#-contribuir-código)
- [Configuración del Entorno](#-configuración-del-entorno)
- [Proceso de Desarrollo](#-proceso-de-desarrollo)
- [Standards de Código](#-standards-de-código)
- [Proceso de Pull Request](#-proceso-de-pull-request)
- [Tipos de Contribuciones](#-tipos-de-contribuciones)

---

## 📜 Código de Conducta

Este proyecto se adhiere a un código de conducta. Al participar, se espera que mantengas este código. Por favor reporta comportamiento inaceptable a través de los [Issues](../../issues).

### Nuestros Estándares

**Comportamiento esperado:**
- ✅ Uso de lenguaje acogedor e inclusivo
- ✅ Respeto a diferentes puntos de vista y experiencias
- ✅ Aceptación de crítica constructiva
- ✅ Enfoque en lo mejor para la comunidad
- ✅ Empatía hacia otros miembros

**Comportamiento inaceptable:**
- ❌ Uso de lenguaje o imágenes sexualizadas
- ❌ Comentarios trolling, insultantes o despectivos
- ❌ Acoso público o privado
- ❌ Publicar información privada de terceros
- ❌ Otras conductas no profesionales

---

## 🚀 ¿Cómo puedo contribuir?

### 🐛 Reportar Bugs

Los bugs se rastrean como [GitHub Issues](../../issues). Antes de crear un bug report:

**Checklist antes de reportar:**
- [ ] Busca en [Issues existentes](../../issues) para evitar duplicados
- [ ] Verifica que estás usando la última versión
- [ ] Asegúrate de que es un bug y no un problema de configuración

**Cómo crear un buen bug report:**

```markdown
## Descripción del Bug
[Descripción clara y concisa]

## Pasos para Reproducir
1. Inicializar servicio con '...'
2. Llamar método '...'
3. Ver error

## Comportamiento Esperado
[Qué esperabas que sucediera]

## Comportamiento Actual
[Qué sucedió en realidad]

## Entorno
- **OS**: [e.g., Ubuntu 22.04]
- **Go Version**: [e.g., 1.25.3]
- **Librería Version**: [e.g., v1.0.0]
- **Valkey/Redis Version**: [e.g., 7.2]

## Logs/Screenshots
[Si aplica, agrega logs o capturas]

## Código para Reproducir
```go
// Código mínimo que reproduce el error
service := anime.NewAnimeFlv()
// ...
```
```

**Ejemplo real:**

> **Título:** SearchAnime retorna error de timeout con búsquedas largas
>
> **Descripción:** Cuando busco animes con nombres muy largos (>50 caracteres), el método `SearchAnime` retorna timeout después de 30 segundos.
>
> **Pasos:** 
> 1. `service.SearchAnime(ctx, "nombre-muy-largo-que-excede-cincuenta-caracteres", 1)`
> 2. Esperar ~30 segundos
> 3. Error: `context deadline exceeded`
>
> **Esperado:** Debería retornar resultados o error específico más rápido
>
> **Entorno:** Go 1.25.3, Ubuntu 22.04, v1.0.0

---

### 💡 Sugerir Features

Las mejoras se rastrean como [GitHub Issues](../../issues) con la etiqueta `enhancement`.

**Checklist antes de sugerir:**
- [ ] Busca en Issues existentes por sugerencias similares
- [ ] Considera si encaja con el alcance del proyecto
- [ ] Piensa en cómo beneficiaría a la mayoría de usuarios

**Cómo crear una buena sugerencia:**

```markdown
## Feature Request: [Título descriptivo]

### Problema a Resolver
[Describe el problema que esta feature resolvería]

### Solución Propuesta
[Describe cómo imaginas la solución]

### Alternativas Consideradas
[Otras soluciones que hayas considerado]

### Casos de Uso
[Ejemplos concretos de cómo se usaría]

### Beneficios
- Beneficio 1
- Beneficio 2

### Posible Implementación
[Opcional: Ideas de cómo implementarlo]
```

**Ejemplo real:**

> **Feature Request:** Soporte para búsqueda de animes por género
>
> **Problema:** Actualmente solo se puede buscar por nombre. Sería útil filtrar por género (acción, drama, etc.)
>
> **Solución:** Nuevo método `SearchByGenre(ctx, genre, page)`
>
> **Casos de Uso:**
> ```go
> // Buscar todos los animes de acción
> results, _ := service.SearchByGenre(ctx, "accion", 1)
> ```
>
> **Beneficios:**
> - Descubrimiento de animes por preferencias
> - Menos carga en servidor (búsquedas específicas)

---

### 💻 Contribuir Código

#### 1. Fork y Clone

```bash
# Fork el repositorio en GitHub, luego:
git clone https://github.com/TU-USUARIO/api-anime.git
cd api-anime

# Agregar upstream
git remote add upstream https://github.com/dst3v3n/api-anime.git
```

#### 2. Crear Rama

```bash
# Actualizar main
git checkout main
git pull upstream main

# Crear rama descriptiva
git checkout -b feature/nombre-descriptivo
# o
git checkout -b fix/descripcion-del-bug
```

**Nomenclatura de ramas:**
- `feature/` - Nuevas funcionalidades
- `fix/` - Corrección de bugs
- `docs/` - Documentación
- `refactor/` - Refactorización
- `test/` - Tests

#### 3. Hacer Cambios

```bash
# Hacer commits pequeños y descriptivos
git add .
git commit -m "feat: agregar búsqueda por género"

# O para bugs
git commit -m "fix: corregir timeout en búsquedas largas"
```

**Formato de commits:**
```
<tipo>: <descripción corta>

[Opcional] Cuerpo del commit explicando:
- Por qué se hizo el cambio
- Qué problema resuelve
- Referencias a issues (#123)
```

**Tipos de commits:**
- `feat:` - Nueva funcionalidad
- `fix:` - Corrección de bug
- `docs:` - Documentación
- `style:` - Formato (no afecta código)
- `refactor:` - Refactorización
- `test:` - Tests
- `chore:` - Mantenimiento

#### 4. Ejecutar Tests

```bash
# Tests unitarios
go test ./test/unit/... -v

# Tests de integración
go test ./test/integration/... -v

# Todos los tests con cobertura
go test ./... -cover

# Verificar que cobertura sea >80%
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

#### 5. Push y Pull Request

```bash
# Push a tu fork
git push origin feature/nombre-descriptivo

# Crear Pull Request en GitHub
```

---

## 🛠️ Configuración del Entorno

### Prerrequisitos

```bash
# Go 1.25.3+
go version

# Valkey/Redis (para tests de integración)
docker run -d -p 6379:6379 valkey/valkey:latest

# Git
git --version
```

### Setup del Proyecto

```bash
# 1. Clone
git clone https://github.com/dst3v3n/api-anime.git
cd api-anime

# 2. Instalar dependencias
go mod download

# 3. Copiar configuración
cp .env.example .env

# 4. Verificar que todo funciona
go test ./... -v
```

### Estructura del Proyecto

```
api-anime/
├── internal/
│   ├── adapters/        # Implementaciones concretas
│   │   ├── cache/       # Adaptador Valkey
│   │   └── scrapers/    # Scrapers (AnimeFlv)
│   ├── domain/          # Lógica de negocio
│   │   ├── dto/         # Data Transfer Objects
│   │   └── services/    # Servicios principales
│   └── ports/           # Interfaces/Contratos
├── config/              # Configuración
├── test/
│   ├── unit/            # Tests unitarios
│   └── integration/     # Tests de integración
├── animeflv.go          # API pública
└── README.md
```

---

## 📏 Standards de Código

### Estilo de Código

```bash
# Formatear código
go fmt ./...

# Linting
go vet ./...

# Golangci-lint (recomendado)
golangci-lint run
```

### Convenciones Go

- ✅ Sigue [Effective Go](https://golang.org/doc/effective_go.html)
- ✅ Usa `gofmt` para formatear
- ✅ Nombres descriptivos (evita abreviaciones)
- ✅ Comentarios en funciones públicas
- ✅ Manejo explícito de errores
- ✅ Interfaces pequeñas y específicas

### Ejemplos de Buen Código

**✅ Correcto:**
```go
// SearchAnime busca animes por nombre con paginación.
// Retorna AnimeResponse con los resultados y el total de páginas.
func (s *AnimeflvService) SearchAnime(ctx context.Context, animeName string, page uint) (dto.AnimeResponse, error) {
    if animeName == "" {
        return dto.AnimeResponse{}, fmt.Errorf("anime name cannot be empty")
    }
    
    // Implementación...
}
```

**❌ Incorrecto:**
```go
// busca animes
func (s *AnimeflvService) search(c context.Context, a string, p uint) (dto.AnimeResponse, error) {
    // sin validación
    // ...
}
```

### Tests

**Cada contribución debe incluir tests:**

```go
// test/unit/animeflv/service_test.go
func TestSearchAnime_ValidInput(t *testing.T) {
    // Arrange
    service := anime.NewAnimeFlv()
    ctx := context.Background()
    
    // Act
    result, err := service.SearchAnime(ctx, "Naruto", 1)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    
    if len(result.Animes) == 0 {
        t.Error("Expected animes, got empty slice")
    }
}

func TestSearchAnime_EmptyName(t *testing.T) {
    service := anime.NewAnimeFlv()
    ctx := context.Background()
    
    _, err := service.SearchAnime(ctx, "", 1)
    
    if err == nil {
        t.Error("Expected error for empty name, got nil")
    }
}
```

### Documentación

```go
// Todas las funciones públicas deben tener comentarios:

// NewAnimeFlv crea una nueva instancia del servicio de AnimeFlv.
// Inicializa automáticamente el scraper y el sistema de caché si está habilitado.
//
// Ejemplo:
//   service := anime.NewAnimeFlv()
//   results, err := service.SearchAnime(ctx, "Naruto", 1)
func NewAnimeFlv() *AnimeflvService {
    // ...
}
```

---

## 🔄 Proceso de Pull Request

### Checklist antes de abrir PR

- [ ] Código formateado con `go fmt ./...`
- [ ] Tests pasan: `go test ./...`
- [ ] Cobertura >80%: `go test ./... -cover`
- [ ] Sin errores de linting: `go vet ./...`
- [ ] Comentarios en funciones públicas
- [ ] README actualizado si es necesario
- [ ] Changelog actualizado (si aplica)
- [ ] Commits descriptivos

### Plantilla de PR

```markdown
## Descripción
[Descripción clara del cambio]

## Tipo de Cambio
- [ ] Bug fix (cambio que corrige un problema)
- [ ] Nueva feature (cambio que agrega funcionalidad)
- [ ] Breaking change (cambio que rompe compatibilidad)
- [ ] Documentación

## ¿Cómo se ha testeado?
[Describe las pruebas realizadas]

## Checklist
- [ ] Mi código sigue el estilo del proyecto
- [ ] He realizado self-review
- [ ] He comentado código complejo
- [ ] He actualizado documentación
- [ ] Mis cambios no generan warnings
- [ ] He agregado tests
- [ ] Tests nuevos y existentes pasan localmente
- [ ] Cobertura >80%

## Screenshots (si aplica)
[Capturas o logs]

## Relacionado
Fixes #[issue number]
```

### Proceso de Review

1. **Automated checks** - CI ejecuta tests y linting
2. **Code review** - Mantenedor revisa el código
3. **Cambios solicitados** - Si se necesitan ajustes
4. **Aprobación** - Una vez que todo esté correcto
5. **Merge** - Se integra a la rama `main`

### Tiempos de Respuesta

- Primera respuesta: 1-3 días hábiles
- Review completo: 3-7 días hábiles
- Merge después de aprobación: 1-2 días

---

## 🎯 Tipos de Contribuciones

### 🐛 Bug Fixes (Bienvenidos siempre)

- Correcciones de errores existentes
- Mejoras de manejo de errores
- Fixes de edge cases

### ✨ Nuevas Features

**Features bienvenidas:**
- Soporte para nuevos sitios de anime
- Mejoras de rendimiento
- Nuevos métodos de búsqueda/filtrado
- Optimizaciones de caché

**Features que requieren discusión:**
- Cambios en API pública (breaking changes)
- Nuevas dependencias
- Cambios arquitectónicos grandes

### 📝 Documentación

- Mejoras en README
- Correcciones de typos
- Ejemplos adicionales
- Traducciones
- Tutoriales

### 🧪 Tests

- Agregar tests faltantes
- Mejorar cobertura
- Tests de integración
- Benchmarks

### 🔧 Refactorización

- Mejoras de código sin cambiar funcionalidad
- Optimizaciones
- Limpieza de código

---

## 💬 ¿Necesitas Ayuda?

- **Documentación**: Lee el [README](README.md) completo
- **Issues**: Busca en [Issues existentes](../../issues)
- **Discussions**: Pregunta en [GitHub Discussions](../../discussions)
- **Email**: Contacta al mantenedor

---

## 🙏 Reconocimientos

¡Todas las contribuciones son valoradas! Los contribuidores serán agregados al README.

### Tipos de Reconocimiento

- 🌟 Mención en CHANGELOG
- 👤 Nombre en lista de contribuidores
- 🏆 Badge especial para contribuciones significativas

---

## 📚 Recursos Adicionales

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)
- [Semantic Versioning](https://semver.org/)

---

**¡Gracias por contribuir a Anime API! 🎉**

*Si tienes dudas sobre el proceso de contribución, no dudes en preguntar abriendo un Issue.*

---

*Última actualización: Diciembre 2024*
