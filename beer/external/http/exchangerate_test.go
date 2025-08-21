// Package http_test contiene pruebas unitarias para el adaptador HTTP de conversión de moneda.
// Utiliza un cliente HTTP simulado para verificar comportamiento de GetExchangeCurrency
// en flujos de éxito y múltiples errores (red, status, lectura, JSON).
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	externalModel "go_trainning/beer-api/beer/external/model"

	"github.com/stretchr/testify/assert"
)

const (
	fakeAmountFloat     = 6.5   // Monto de ejemplo a convertir
	fakeCurrencyFromStr = "EUR" // Moneda origen para pruebas
	fakeCurrencyToStr   = "USD" // Moneda destino para pruebas
)

var (
	fakeQuoteFloat  = 1.15553           // Tipo de cambio simulado
	fakeResultFloat = 7.510945          // Resultado esperado de la conversión
	fakeTimeUnix    = time.Now().Unix() // Timestamp simulado
)

// badReader simula un Body con fallo al leer.
type badReader struct{}

func (br badReader) Read(p []byte) (int, error) { return 0, errors.New("read failure") }
func (br badReader) Close() error               { return nil }

// RoundTripFunc permite definir un Transport HTTP personalizado para pruebas.
type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

// newTestClient retorna un *http.Client que usa el RoundTripFunc dado como Transport.
func newTestClient(fn RoundTripFunc) *http.Client { return &http.Client{Transport: fn} }

// sampleSuccessBody reconstruye un CurrencyConversionResponse válido en bytes JSON.
func sampleSuccessBody() []byte {
	resp := externalModel.CurrencyConversionResponse{
		Success: true,
		Query:   externalModel.CurrencyQuery{From: fakeCurrencyFromStr, To: fakeCurrencyToStr, Amount: fakeAmountFloat},
		Info:    externalModel.InfoResponse{Timestamp: fakeTimeUnix},
		Result:  fakeResultFloat,
	}
	b, _ := json.Marshal(resp)
	return b
}

// Test_clientExchanGerate_GetExchangeCurrency cubre los siguientes escenarios:
// - request OK con JSON válido
// - error en ejecución HTTP
// - status code != 200
// - error al leer body
// - error de deserialización JSON
func Test_clientExchanGerate_GetExchangeCurrency(t *testing.T) {
	const fakeKey = "TESTKEY123"
	// Configurar API_KEY_EXCHANGERATE para construir URL
	err := os.Setenv("API_KEY_EXCHANGERATE", fakeKey)
	assert.NoError(t, err)
	defer os.Unsetenv("API_KEY_EXCHANGERATE")

	// baseURLCheck valida que los parámetros query sean correctos	notm
	baseURLCheck := func(t *testing.T, req *http.Request) {
		u := req.URL.Query()
		assert.Equal(t, fakeKey, u.Get("access_key"))
		assert.Equal(t, fakeCurrencyToStr, u.Get("from"))
		assert.Equal(t, fakeCurrencyFromStr, u.Get("to"))
	}

	t.Run("success", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) (*http.Response, error) {
			baseURLCheck(t, req)
			return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(sampleSuccessBody()))},
				nil
		})

		api := NewClientExchanGerate(client)
		gotExchange, errExchangeCurrency := api.GetExchangeCurrency(fakeCurrencyFromStr, fakeCurrencyToStr, fakeAmountFloat)

		assert.NoError(t, errExchangeCurrency)
		assert.Equal(t, fakeResultFloat, gotExchange.Result)
		assert.Equal(t, fakeTimeUnix, gotExchange.Info.Timestamp)
	})

	t.Run("HTTP request error", func(t *testing.T) {
		client := newTestClient(func(_ *http.Request) (*http.Response, error) {
			return nil, errors.New("network failure")
		})

		api := NewClientExchanGerate(client)
		_, errExchangeCurrency := api.GetExchangeCurrency(fakeCurrencyFromStr, fakeCurrencyToStr, fakeAmountFloat)

		assert.Error(t, errExchangeCurrency)
		assert.Contains(t, errExchangeCurrency.Error(), "network failure")
	})

	t.Run("status != 200", func(t *testing.T) {
		client := newTestClient(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte{}))},
				nil
		})

		api := NewClientExchanGerate(client)
		_, errExchangeCurrency := api.GetExchangeCurrency(fakeCurrencyFromStr, fakeCurrencyToStr, fakeAmountFloat)

		exp := fmt.Sprintf("status code error: %d %s", http.StatusInternalServerError, "500 Internal Server Error")
		assert.Error(t, errExchangeCurrency)
		assert.EqualError(t, errExchangeCurrency, exp)
	})

	t.Run("error reading body", func(t *testing.T) {
		client := newTestClient(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
					StatusCode: http.StatusOK,
					Body:       badReader{}},
				nil
		})

		api := NewClientExchanGerate(client)
		_, errExchangeCurrency := api.GetExchangeCurrency(fakeCurrencyFromStr, fakeCurrencyToStr, fakeAmountFloat)

		assert.Error(t, errExchangeCurrency)
		assert.EqualError(t, errExchangeCurrency, "read failure")
	})

	t.Run("error JSON Unmarshal", func(t *testing.T) {
		client := newTestClient(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("not a json")))},
				nil
		})

		api := NewClientExchanGerate(client)
		_, errExchangeCurrency := api.GetExchangeCurrency(fakeCurrencyFromStr, fakeCurrencyToStr, fakeAmountFloat)

		var syntaxErr *json.SyntaxError
		assert.Error(t, errExchangeCurrency)
		assert.ErrorAs(t, errExchangeCurrency, &syntaxErr)
	})
}
