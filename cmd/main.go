// Package main es el punto de entrada de la aplicación Beer API.
// Se encarga de:
//  1. Cargar variables de entorno desde .env (opcional).
//  2. Construir y configurar el contenedor de inyección de dependencias.
//  3. Inicializar el logger según configuración (modo debug opcional).
//  4. Resolver e invocar componentes: servidor HTTP, router y base de datos.
//  5. Aplicar migraciones y arrancar el servidor.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	// Contenedor DI
	"go_trainning/beer-api/configs/generals/injector"
	// Configuración de rutas
	"go_trainning/beer-api/configs/generals/router"
	// Conexión y migraciones BD
	"go_trainning/beer-api/configs/storage"
	// Utilidades (logger)
	"go_trainning/beer-api/pkg/kit"

	// Carga variables de entorno desde archivo .env
	"github.com/joho/godotenv"
	// Framework web Echo
	"github.com/labstack/echo/v4"
	// Logger zerolog
	"github.com/rs/zerolog/log"
)

// main es la función principal que inicializa y arranca la aplicación.
func main() {
	// 1. Cargar configuración de entorno desde archivo .env
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("⚠️  No se encontró archivo .env, usaré variables de entorno del sistema")
	}

	// 2. Construir el contenedor de dependencias
	container := injector.BuildContainer()

	// 3. Determinar si el logger debe estar en modo debug según variable de entorno
	boolVal, errBool := strconv.ParseBool(os.Getenv("LOGGER_DEBUG"))
	if errBool != nil {
		panic(fmt.Errorf("LOGGER_DEBUG must be set to true or false: %w", errBool))
	}

	// Definir flag de debug para permitir override al ejecutar
	debug := flag.Bool("debug", boolVal, "sets log level to debug")
	flag.Parse()

	// 4. Inicializar el logger global con nombre de aplicación y nivel
	kit.InitLogger("beer-api", *debug)

	// Resolver e invocar la función que arranca la aplicación
	err := container.Invoke(func(server *echo.Echo, route *router.Router, db *sql.DB) {
		// Construir dirección de escucha usando puerto de entorno
		address := fmt.Sprintf("%s:%s", "0.0.0.0", os.Getenv("API_PORT"))

		// Inicializar rutas en el servidor Echo
		route.Init()

		// Aplicar migraciones de base de datps antes de servir
		err := storage.VersionedDB(db)
		if err != nil {
			panic(fmt.Errorf("error executing migrations: %v", err))
		}

		// Arrancar servidor HTTP y registrar fatal si falla
		server.Logger.Fatal(server.Start(address))
	})

	// Manejo de errores al invocar contenedor DI
	if err != nil {
		panic(fmt.Sprintf("Error initializing aplication: %v", err))
	}

	// 6. Asegurar cierre de conexión a BD al finalizar
	defer storage.PostgresCloseConnection()
}