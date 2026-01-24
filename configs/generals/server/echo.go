// Package server contiene la configuración básica del servidor HTTP
// utilizando el framework Echo. Aquí se inicializan los middlewares
// comunes que debe usar toda la aplicación (CORS, logging, manejo de errores).
package server

import (
	// Middleware CORS personalizado
	"go_trainning/beer-api/pkg/middleware"

	// Echo: framework web minimalista de alto rendimiento
	"github.com/labstack/echo/v4"
	// Middlewares integrados (RequestID, Logger, Recover)
	middlewareEcho "github.com/labstack/echo/v4/middleware"
)

// NewServer crea y devuelve una instancia de *echo.Echo ya configurada.
// Pasos de configuración:
//  1. Crear nueva instancia de Echo.
//  2. Aplicar middleware de CORS usando la configuración predeterminada.
//  3. Agregar RequestID para trazar cada petición con un UUID único.
//  4. Registrar Logger para capturar metodo, ruta, código de estado y latencia.
//  5. Registrar Recover para interceptar panics y devolver HTTP 500 sin crash.
//
// Retorna:
//   - *echo.Echo: servidor listo para arrancar con Start().
func NewServer() *echo.Echo {
	// 1. Crear nueva instancia de Echo
	server := echo.New()

	// 2. Middleware CORS personalizado (orígenes, métodos, headers configurables)
	server.Use(middleware.CORSMiddleware(middleware.DefaultCORSConfig))

	// 3. Middleware RequestID: añade un identificador único por petición
	server.Use(middlewareEcho.RequestID())

	// 4. Logger: registra detalles de cada petición
	server.Use(middlewareEcho.Logger())

	// 5. Recover: captura panics y responde 500 sin detener el servidor
	server.Use(middlewareEcho.Recover())

	return server
}