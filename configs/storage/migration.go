// Package storage proporciona utilidades para la gestión de migraciones
// de esquema de la base de datos utilizando la librería golang-migrate.
package storage

import (
	// database/sql define la interfaz genérica para operaciones SQL.
	"database/sql"
	// fmt para formatear mensajes de error.
	"fmt"
	// log para registrar mensajes críticos en stdout/stderr.
	"log"
	// os para leer variables de entorno de configuración.
	"os"
	// strings para operaciones sobre cadenas (p.ej. detección de substrings).
	"strings"

	// migrate es el paquete principal de migraciones.

	// database contiene interfaces y utilidades para drivers de BD.
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"

	// postgres contiene la implementación del driver de migración para PostgreSQL.
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// source/file registra el origen de migraciones desde archivos.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// Registro del driver de PostgreSQL para ping y sql.Open.
	_ "github.com/lib/pq"
)

// VersionedDB ejecuta las migraciones “UP” sobre la base de datos provista.
//   1. Verifica conectividad con db.Ping().
//   2. Construye el driver de migración según DB_DRIVER.
//   3. Llama a migrationUp para aplicar los scripts.
//   4. Ignora el error “no change” si no hay cambios pendientes.
//
// Variables de entorno requeridas:
//   - DB_DRIVER    : driver de la BD (e.g., "postgres")
//   - DB_NAME      : nombre lógico usado en migrate NewWithDatabaseInstance
//   - SCRIPTS_PATH : ruta al directorio de scripts de migración
func VersionedDB(db *sql.DB) error {
	// Asegurar que la conexión está viva
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	// Configurar el driver de migración
	var instanceConfig database.Driver
	var errConfig error
	switch os.Getenv("DB_DRIVER") {
		case "postgres":
			// Crea una instancia de driver Postgres para migraciones
			instanceConfig, errConfig = postgres.WithInstance(db, &postgres.Config{})
		default:
			log.Fatalf("unsupported DB_DRIVER: %s", os.Getenv("DB_DRIVER"))
	}

	if errConfig != nil {
		log.Fatal("error configuring migration driver:", errConfig)
	}

	// Ejecutar las migraciones “Up”
	version, err := migrationUp(instanceConfig)
	if err != nil {
		// Si el error contiene “no change”, ignorarlo; en otro caso, envolverlo
		if strings.Contains(err.Error(), "no change") {
			return nil
		}
		return fmt.Errorf("migration failed at version %d: %w", version, err)
	}

	return nil
}

// migrationUp aplica hacia arriba (“Up”) todos los scripts de migración.
// Devuelve la versión final aplicada y cualquier error que no sea “no change”.
func migrationUp(instanceConfig database.Driver) (int, error) {
	// Leer configuración de entorno
	pathScripts := os.Getenv("SCRIPTS_PATH")
	dbName := os.Getenv("DB_NAME")

	// Crear el objeto migrate que apunta a los scripts y al driver de BD
	migrator, err := migrate.NewWithDatabaseInstance(
		pathScripts,
		dbName,
		instanceConfig,
	)
	if err != nil {
		log.Fatalf("failed to initialize migrator: %s", err)
	}

	// Ejecutar todas las migraciones pendientes
	err = migrator.Up()
	if err != nil {
		// Si no es “no change”, forzar rollback una versión y reportar
		if !strings.Contains(err.Error(), "no change") {
			ver, _, verErr := migrator.Version()
			if verErr != nil {
				log.Fatal("failed to get migration version:", verErr)
			}
			// Retroceder al estado previo a la última migración fallida
			if forceErr := migrator.Force(int(ver) - 1); forceErr != nil {
				log.Fatalf("failed to force migration rollback: %s", forceErr)
			}
			return int(ver), err
		}
	}

	// Obtener la versión actual después de Up() exitoso o “no change”
	finalVer, _, _ := migrator.Version()
	return int(finalVer), nil
}