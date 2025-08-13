package main

import (
	// router agrupa la configuración de rutas y middlewares sobre Echo.
	"go_trainning/beer-api/configs/generals/router"
	// server proporciona la instancia de Echo con middlewares globales (CORS).
	"go_trainning/beer-api/configs/generals/server"
)

// main ejecuta el flujo de arranque de la aplicación:
//   1. Crear la instancia de Echo configurada con CORS.
//   2. Registrar rutas y middlewares usando el módulo router.
//   3. Iniciar el servidor en el puerto 8888 y, en caso de error, terminar la ejecución.
func main() {
	// Inicializa el servidor con configuración básica (CORS).
	serverEcho := server.NewServer()

	// Registra rutas y middlewares adicionales (RequestID, Logger, Recover, endpoints).
	router.Init(serverEcho)

	// Arranca el servidor en el puerto 8888; registra fatal si ocurre un error.
	serverEcho.Logger.Fatal(serverEcho.Start(":8888"))
}