// Package injector configura el contenedor de dependencias usando dig (Uber Dig).
// Define y provee todos los componentes de la aplicación, conectando su constructor
// con sus dependencias para crear un grafo de inyección.
package injector

import (
	"fmt"
	"net/http"
	"time"

	// Módulos de aplicación:
	externalHttp "go_trainning/beer-api/beer/external/http"
	"go_trainning/beer-api/beer/handler"
	"go_trainning/beer-api/beer/repository"
	"go_trainning/beer-api/beer/service"
	"go_trainning/beer-api/configs/generals/router"
	"go_trainning/beer-api/configs/generals/server"
	"go_trainning/beer-api/configs/storage"

	// dig es el contenedor de inyección de dependencias de Uber.
	"go.uber.org/dig"
)

// Container mantiene la referencia global al contenedor de dependencias.
var Container *dig.Container

// BuildContainer construye y configura el grafo de dependencias
// Registra (Provide) cada constructor con sus dependencias en el orden:
//  1. storage.PostgresConnection     -> *sql.DB (conexión a Postgres)
//  2. storage.VersionedDB            -> error (ejecuta migraciones)
//  3. router.NewRouter               -> *router.Router
//  4. server.NewServer               -> *echo.Echo
//  5. handler.NewBeerHandler         -> handler.BeerHandler
//  6. service.NewBeerService         -> service.BeerService
//  7. repository.NewBeerRepository   -> repository.BeerRepository
//  8. func() *http.Client            -> *http.Client (configurado con timeout)
//  9. externalHttp.NewClientExchanGerate -> interfaces.CurrencyLayer
//
// Retorna el contenedor listo para invocaciones.
func BuildContainer() *dig.Container {
	Container = dig.New()

	// 1. Conexión a la base de datos
	checkError(Container.Provide(storage.PostgresConnection))
	// 2. Migraciones (opcional) - descomentar si se usan
	// checkError(Container.Provide(storage.VersionedDB))

	// 3. Infraestructura HTTP y rutas
	checkError(Container.Provide(router.NewRouter))
	checkError(Container.Provide(server.NewServer))
	
	// 4. Capa de Handlers
	checkError(Container.Provide(handler.NewBeerHandler))

	// 5. Capa de Servicios
	checkError(Container.Provide(service.NewBeerService))

	// 6. Capa de Repositorio
	checkError(Container.Provide(repository.NewBeerRepository))

	// 7. Cliente HTTP
	checkError(Container.Provide(func() *http.Client {
		return &http.Client{Timeout: 15 * time.Second}
	}))

	// 8. Cliente externo para conversión de moneda
	checkError(Container.Provide(externalHttp.NewClientExchanGerate))

	return Container
}

// checkError paniquea si ocurre un error al registrar un constructor.
func checkError(err error) {
	if err != nil {
		panic(fmt.Sprintf("error injecting dependencies: %v", err))
	}
}