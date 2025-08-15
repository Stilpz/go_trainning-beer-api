// Package router se encarga de exponer y organizar los endpoints HTTP de la aplicación.
// Define rutas, handlers y realiza el registro de logs de cada ruta.
package router

import (
	"net/http"

	// Echo es el framework web utilizado para definir rutas y handlers.
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	// zerolog para registro estructurado de las rutas cargadas.
	"github.com/rs/zerolog/log"

	// handler contiene las funciones que atienden las solicitudes HTTP.
	"go_trainning/beer-api/beer/handler"
)

// Router agrupa la instancia de Echo y los handlers asociados.
// Se encarga de inicializar rutas y middlewares específicos de la aplicación.
type Router struct {
	server      *echo.Echo          // Instancia de Echo con middlewares globales
	beerHandler handler.BeerHandler // Handler que delega la lógica de BeerService
}

// healthCheckResponse representa la respuesta JSON del endpoint /health.
type healthCheckResponse struct {
	Status string `json:"status"` // Indica el estado de salud del servicio
}

// NewRouter construye un nuevo Router con la instancia de Echo y el BeerHandler.
// Parámetros:
//   - server: *echo.Echo con middlewares ya configurados (CORS, logger, recover, etc.)
//   - beerHandler: implementación de handler.BeerHandler para endpoints de cerveza
//
// Retorna:
//   - *Router: router listo para inicializar rutas.
func NewRouter(
	server *echo.Echo,
	beerHandler handler.BeerHandler,
) *Router {
	return &Router{
		server:      server,
		beerHandler: beerHandler,
	}
}

// Init registra las rutas HTTP sobre la instancia de Echo asociada.
// Rutas configuradas:
//
//	GET /health -> healthCheckHandler
//	GET / -> Listar todas las cervezas
//	GET /:beerID -> Obtener detalles de una cerveza
//	GET /:beerID/box-price -> Calcular precio de caja de cerveza
//	POST / -> Crear nueva cerveza
//
// Además, itera sobre todas las rutas registradas y escribe un log con metodo y path.
func (r *Router) Init() {
	apiGroup := r.server.Group("")

	// Health check
	apiGroup.GET("/health", healthCheckHandler)
	apiGroup.GET("/docs/*", echoSwagger.WrapHandler)

	// Endpoints de Beer
	apiGroup.GET("/", r.beerHandler.GetAllBeersHandler)
	apiGroup.GET("/:beerID", r.beerHandler.GetOneHandler)
	apiGroup.GET("/:beerID/box-price", r.beerHandler.GetOneBoxPriceHandler)
	apiGroup.POST("/", r.beerHandler.CreateHandler)

	// Loguear cada ruta registrada para monitoreo
	for _, route := range r.server.Routes() {
		log.Info().Msgf("[%s] %s", route.Method, route.Path)
	}
}

// healthCheckHandler responde con el estado de salud del servicio.
// Retorna HTTP 200 con JSON {"status":"ok"}.
func healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, healthCheckResponse{Status: "ok"})
}