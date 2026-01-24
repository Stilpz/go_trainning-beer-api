// Package http provee adaptadores HTTP para comunicarse con servicios externos.
// En este caso, implementa la consulta de tasas de cambio a través de exchangerate.host.
package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// externalModel contiene la estructura de la respuesta de conversión de moneda.
	externalModel "go_trainning/beer-api/beer/external/model"
	// interfaces define el contrato CurrencyLayer que debe implementar este cliente.
	"go_trainning/beer-api/beer/interfaces"
)

// clientExchanGerate es la implementación concreta de interfaces.CurrencyLayer
// que utiliza un *http.Client para realizar las peticiones HTTP.
type clientExchanGerate struct {
	clientHttp *http.Client
}

// NewClientExchanGerate construye un nuevo adaptador CurrencyLayer.
// Parámetros:
//   - httpClient: instancia de *http.Client preconfigurada (timeout, transporte, etc.).
//
// Retorna:
//   - interfaces.CurrencyLayer: cliente listo para invocar GetCurrency.
func NewClientExchanGerate(httpClient *http.Client) interfaces.CurrencyLayer {
	return &clientExchanGerate{
		clientHttp: httpClient,
	}
}

// GetExchangeCurrency realiza una llamada a la API de exchangerate.host para convertir un monto
// de currencyPay a currencyBeer en la fecha actual.
// Construye la URL usando la API key leída de la variable de entorno API_KEY_EXCHANGERATE.
//
// Parámetros:
//   - currencyPay: código ISO de la moneda de origen (e.g., "USD").
//   - currencyBeer: código ISO de la moneda destino (e.g., "EUR").
//   - amountBeer: monto en currencyPay a convertir.
//
// Retorna:
//   - externalModel.CurrencyConversionResponse: datos de la conversión (tasa, monto convertido, fecha).
//   - error: en caso de fallo de conexión, código HTTP distinto de 200 o error de decodificación JSON.
func (c *clientExchanGerate) GetExchangeCurrency(currencyPay, currencyBeer string, amountBeer float64) (
	externalModel.CurrencyConversionResponse, error) {
	// Valor por defecto en caso de error
	var emptyResp externalModel.CurrencyConversionResponse

	// Leer API key y formatear la fecha actual (YYYY-MM-DD)
	accessKey := os.Getenv("API_KEY_EXCHANGERATE")
	currentDate := time.Now().Format("2006-01-02")

	// Construir y ejecutar la petición HTTP GET
	url := fmt.Sprintf("https://api.exchangerate.host/convert?access_key=%s&from=%s&to=%s&amount=%f&date=%s",
		accessKey, currencyBeer, currencyPay, amountBeer, currentDate)
	response, err := c.clientHttp.Get(url)
	if err != nil {
		log.Println("error client http execute :", err)
		return emptyResp, err
	}
	// Asegurar cierre de Body al terminar
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Println("error closing response body:", closeErr)
		}
	}()

	// Verificar código de estado HTTP
	if response.StatusCode != http.StatusOK {
		return emptyResp, fmt.Errorf("status code error: %d %s", response.StatusCode, response.Status)
	}

	// Leer cuerpo de la respuesta
	bodyBytes, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		log.Println("error read body :", errRead)
		return emptyResp, errRead
	}

	// Decodificar JSON en la estructura de respuesta
	var conversion externalModel.CurrencyConversionResponse
	if errUnmarshal := json.Unmarshal(bodyBytes, &conversion); errUnmarshal != nil {
		log.Println("error json unmarshal body:", errUnmarshal)
		return emptyResp, errUnmarshal
	}

	return conversion, nil
}