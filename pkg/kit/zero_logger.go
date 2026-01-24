// Package kit provee utilidades transversales para la aplicación.
// En particular, zero_logger configura el logger global usando zerolog
// en un formato de consola legible con colores por nivel.
package kit

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Constantes de formato de color ANSI para los niveles de log
const (
	InfoColor    	= "\033[1;32m%s\033[0m"    // Verde brillante
	WarningColor	= "\033[1;33m%s\033[0m"    // Amarillo brillante
	ErrorColor   	= "\033[1;31m%s\033[0m"    // Rojo brillante
)

// InitLogger inicializa el logger global de la aplicación.
// Parámetros:
//   - appName: nombre de la aplicación, añadido como campo "app" en cada log.
//   - debug: si es true, habilita el nivel de log Debug; de lo contrario, Info.
//
// Comportamiento:
//  1. Configura zerolog.ConsoleWriter para salida en stderr con timestamp RFC3339.
//  2. Personaliza FormatLevel para envolver el nivel en colores ANSI.
//  3. Define FormatMessage y FormatFieldName para mejorar legibilidad.
//  4. Asigna log.Logger al nuevo logger configurado, incluyendo timestamp, caller y campo "app".
//  5. Establece el nivel global: Debug si debug=true, Info por defecto.
//  6. Si el nivel Debug está habilitado, escribe un mensaje de confirmación.
func InitLogger(appName string, debug bool) {
	// Configurar salida en consola con colores y timestamps
	output := zerolog.ConsoleWriter{
		Out: 	  os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:  false,
	}

	// Formato del nivel colores dependiendo del valor
	output.FormatLevel = func(i interface{}) string {
		level := strings.ToLower(i.(string))
		var colored string
		switch level {
		case "warn":
			colored = fmt.Sprintf(WarningColor, strings.ToUpper(level))
		case "error":
			colored = fmt.Sprintf(ErrorColor, strings.ToUpper(level))
		default:
			colored = fmt.Sprintf(InfoColor, strings.ToUpper(level))
		}
		return fmt.Sprintf("[%v]", colored)
	}

	// Formateo del mensaje
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s |", i)
	}

	// Formateo del nombre del campo
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}

	// Construir logger con timestamp, caller y campo "app"
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger().With().Str("app", appName).Logger()

	// Nivel global por defecto
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Si debug, elevar nivel para mostrar Debug logs
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Registrar mensaje si Debug está habilitado
	if log.Debug().Enabled() {
		log.Debug().Msg("Debug mode enabled")
	}
}