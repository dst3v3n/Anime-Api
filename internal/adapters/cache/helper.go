// Package cache contiene funciones auxiliares para serialización y deserialización
// de datos utilizados en el sistema de caché del sistema.
package cache

import "encoding/json"

// deserialize convierte una cadena JSON en una estructura de destino.
// Si la cadena está vacía, retorna sin error (valor por defecto).
// Utiliza la función json.Unmarshal para parsear el contenido JSON.
func deserialize(data string, dest interface{}) error {
	if data == "" {
		return nil // o un error personalizado si prefieres
	}
	return json.Unmarshal([]byte(data), dest)
}
