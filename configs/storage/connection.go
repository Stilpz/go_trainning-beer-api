// Package storage proporciona la configuración y gestión de la conexión
// a la base de datos PostgreSQL. Implementa un patrón singleton para
// garantizar una única instancia de *sql.DB en toda la aplicación.
package storage

import (
	// database/sql ofrece una interfaz genérica para bases de datos SQL.
	"database/sql"
	// fmt se utiliza para construir cadenas y formatear errores.
	"fmt"
	// net/url permite escapar caracteres en la contraseña al formar la URL.
	"net/url"
	// os para leer variables de entorno con la configuración del DB.
	"os"
	// sync para garantizar inicialización de singleton sólo una vez.
	"sync"

	// Registro del driver de PostgreSQL.
	_ "github.com/lib/pq"
	// zerolog para registro estructurado de eventos.
	"github.com/rs/zerolog/log"
)

var (
	// once asegura que la conexión se inicialice sólo una vez.
	once sync.Once
	// instance mantiene la única instancia de *sql.DB utilizada
	instance   *sql.DB
)

// PostgresConnection devuelve la instancia singleton de *sql.DB
// La primera vez que se invoca:
// 1. Llama internamente a getConnection para crear la conexión
// 2. Si ocurre un error, registra un Fatal y detiene la aplicación.
// 3. Almacena la conexión en la variable singleton `instance`
//
// Retorna la misma instancia en llamadas posteriores.
func PostgresConnection() *sql.DB {
	db, err := getConnection()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Error connecting to the database: %v", err))
	}

	// Inicializar singleton la primera vez
	once.Do(func() {
		instance = db
	})
	return instance
}

// getConnection lee la configuración de la base de datos desde variables
// de entorno, construye la cadena de conexión, abre y verifica la conexión.
// Variables de entorno utilizadas:
//   - DB_HOST: host de la base de datos
//   - DB_PORT: puerto de conexión
//   - DB_USER: usuario de la base de datos
//   - DB_PASSWORD: contraseña (se escapa para URL)
//   - DB_NAME: nombre de la base de datos
//   - DB_DRIVER: driver SQL (ej. "postgres")
//   - DB_SSL_MODE: modo SSL (ej. "disable", "require")
//
// Retorna *sql.DB si la conexión es exitosa, o es un error descriptivo
func getConnection() (*sql.DB, error) {
	// Leer variables de entorno
	DbHost := os.Getenv("DB_HOST")
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")
	DbSslMode := os.Getenv("DB_SSL_MODE")

	// Construir la URL de conexión, escapando la contraseña
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		DbUser, url.QueryEscape(DbPassword), DbHost, DbPort, DbName, DbSslMode,
	)

	// Abrir conexión
	db, err := sql.Open(DbDriver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err.Error())
	}

	// Validar instancia no nula
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Verificar conexión
	if err = db.Ping(); err != nil {
		// Intentar cerrar en caso de fallo
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to the database %s", DbName)

	return db, nil
}

// PostgresCloseConnection cierra la conexión singleton de PostgreSQL si existe.
// Registra Fatal en caso de error al cerrar.
func PostgresCloseConnection() {
	if instance != nil {
		if err := instance.Close(); err != nil {
			log.Fatal().Msg(fmt.Sprintf("Error closing the database: %v", err))
		}
	}
}