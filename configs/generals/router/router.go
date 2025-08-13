// Package router se encarga de registrar las rutas HTTP y middlewares
// sobre la instancia de Echo provista por el servidor.
package router

import (
	// Echo es el framework web que utilizamos para definir rutas y handlers.
	"github.com/labstack/echo/v4"
	// middlewareEcho provee middlewares integrados como RequestID, Logger y Recover.
	middlewareEcho "github.com/labstack/echo/v4/middleware"
	// zerolog se usa para registrar información estructurada de las rutas cargadas.
	"github.com/rs/zerolog/log"
	// handler contiene las funciones que atienden las solicitudes HTTP.
	"go_trainning/beer-api/beer/handler"
)

// Init recibe una instancia de *echo.Echo y realiza:
//   1. Registro de middlewares globales: RequestID, Logger, Recover.
//   2. Creación de un grupo de rutas (apiGroup) para organizar endpoints.
//   3. Definición de la ruta GET /greeting y su handler correspondiente.
//   4. Logueo de cada ruta registrada para facilitar el monitoreo.
func Init(e *echo.Echo) {
	// Agrega un identificador único (UUID) a cada petición entrante
	e.Use(middlewareEcho.RequestID())
	// Registra la información de cada petición (método, path, status, tiempo)
	e.Use(middlewareEcho.Logger())
	// Captura panics en los handlers y responde con HTTP 500 sin detener el servidor.
	e.Use(middlewareEcho.Recover())

	// Define un grupo de rutas bajo el path raíz (“”).
	apiGroup := e.Group("")

	// Endpoint de prueba para saludar; devuelve un mensaje fijo.
	apiGroup.GET("/greeting", handler.Greeting)

	// Itera sobre todas las rutas registradas y las imprime en el log.
	for _, r := range e.Routes() {
		log.Info().Msgf("[%s] %s", r.Method, r.Path)
	}
}