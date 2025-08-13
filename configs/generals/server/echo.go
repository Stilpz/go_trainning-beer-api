// Package server contiene la configuración básica del servidor HTTP
// utilizando el framework Echo. Aquí se inicializan los middlewares
// comunes que debe usar toda la aplicación.
package server

import (
	// Echo es el framework web de alto rendimiento y minimalist
	"github.com/labstack/echo/v4"
	// middleware incluye componentes como CORS, Logger, Recover, etc.
	"github.com/labstack/echo/v4/middleware"
)

// NewServer crea y devuelve una instancia de *echo.Echo ya configurada.
// Se encarga de:
//   1. Inicializar el servidor Echo.
//   2. Registrar los middlewares globales necesarios (actualmente CORS).
//   3. Retornar la instancia lista para arrancar.
func NewServer() *echo.Echo {
	// Crear una nueva instancia de Echo
	server := echo.New()

	// Registrar middleware de CORS para permitir solicitudes cross-origin.
	// Por defecto, permite todos los orígenes y métodos. Puede personalizarse
	// añadiendo opciones a middleware.CORSWithConfig(...)
	server.Use(middleware.CORS())

	// Devolver la instancia configurada
	return server
}