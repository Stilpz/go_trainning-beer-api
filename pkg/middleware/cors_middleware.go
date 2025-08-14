// Package middleware contiene componentes de middleware reutilizables para Echo.
// En particular, este archivo define un middleware de CORS configurable.
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// CORSConfig contiene la configuración de CORS para el middleware.
//
//	AllowOrigins: lista de orígenes permitidos (ej. "*" para todos).
//	AllowMethods: métodos HTTP permitidos (GET, POST, etc.).
//	AllowHeaders: headers permitidos en la petición.
//	ExposeHeaders: headers que se exponen al cliente en la respuesta.
//	AllowCredentials: si se permite el envío de credenciales (cookies, auth).
//	MaxAge: tiempo máximo (duración) que el preflight puede cachearse.
type CORSConfig struct {
	AllowOrigins	 []string	// Orígenes permitidos
	AllowMethods	 []string	// Métodos HTTP permitidos
	AllowHeaders	 []string	// Headers permitidos en la petición
	ExposeHeaders	 []string	// Headers expuestos en la respuesta
	AllowCredentials bool		// Permitir envío de cookies/credenciales
	MaxAge		 time.Duration	// Duración del cache para preflight
}

// DefaultCORSConfig es una configuración inicial razonable para CORS.
// Se puede utilizar como base y ajustar según el entorno de despliegue.
var DefaultCORSConfig = CORSConfig{
	AllowOrigins:     []string{"http://localhost"}, // Ajustar a dominios reales
	AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderAccept},
	ExposeHeaders:    []string{"Link"},
	AllowCredentials: true,
	MaxAge:          300 * time.Second,
}

// CORSMiddleware crea un middleware Echo que aplica las reglas CORS definidas en config.
// Utiliza los valores de CORSConfig para establecer los headers adecuados en la respuesta.
func CORSMiddleware(config CORSConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			origin := req.Header.Get(echo.HeaderOrigin)

			// Aplicar CORS solo si hay Origin y está permitido
			if origin != "" && isOriginAllowed(origin, config.AllowOrigins) {
				// Reflejar el origen solicitado
				res.Header().Set(echo.HeaderAccessControlAllowOrigin, origin)
				res.Header().Set(echo.HeaderVary, echo.HeaderOrigin)

				// Métodos y headers permitidos
				res.Header().Set(echo.HeaderAccessControlAllowMethods, strings.Join(config.AllowMethods, ","))
				res.Header().Set(echo.HeaderAccessControlAllowHeaders, strings.Join(config.AllowHeaders, ","))
				res.Header().Set(echo.HeaderAccessControlAllowCredentials, strconv.FormatBool(config.AllowCredentials))

				// Headers expuestos si los hay
				if len(config.ExposeHeaders) > 0 {
					res.Header().Set(echo.HeaderAccessControlExposeHeaders, strings.Join(config.ExposeHeaders, ","))
				}

				// Cache max-age para la petición preflight
				res.Header().Set(echo.HeaderAccessControlMaxAge, strconv.Itoa(int(config.MaxAge.Seconds())))
			}

			// Si es una petición OPTIONS (preflight), responder sin contenido (con 204 No Content)
			if req.Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			// Continuar con el siguiente handler en la cadena
			return next(c)
		}
	}
}

// isOriginAllowed verifica si el origen de la petición está en la lista de permitidos.
// Retorna true si allowList contiene "*" o conincide (ignora mayúsculas) con el origen.
func isOriginAllowed(origin string, allowList []string) bool {
	for _, ao := range allowList {
		if ao == "*" || strings.EqualFold(ao, origin) {
			return true
		}
	}
	return false
}
