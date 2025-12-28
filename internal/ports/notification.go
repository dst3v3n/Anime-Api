// Package ports define las interfaces (puertos) que establecen contratos entre
// las diferentes capas de la aplicación siguiendo la arquitectura hexagonal.
//
// notification.go define NotificationPort, la interfaz para implementar
// sistemas de notificación en la aplicación. Esta interfaz está disponible
// para futuras extensiones de funcionalidad de notificación (emails, webhooks, etc.)
// sin acoplar la lógica de negocio a una implementación específica.
package ports

// NotificationPort define el contrato para cualquier sistema de notificación que se implemente.
// Actualmente es una interfaz vacía pero proporciona un punto de extensión para agregar
// métodos de notificación (envío de emails, alertas, webhooks, etc.) en el futuro.
type NotificationPort interface{}
