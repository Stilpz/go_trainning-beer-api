// Package handler_test contiene pruebas unitarias para los handlers HTTP de Beer.
// Utiliza Echo para crear contextos de petición, testify para aserciones y mocks
// para simular comportamientos de BeerService.
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_mocksService "go_trainning/beer-api/beer/mocks/interfaces"
	"go_trainning/beer-api/beer/model"
	"go_trainning/beer-api/beer/repository"
	"go_trainning/beer-api/beer/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	fakeBeerIdUint  uint   = 1     // ID de cerveza de ejemplo
	fakeQuantityInt int    = 1     // Cantidad de ejemplo
	fakeCurrencyStr string = "USD" // Moneda de ejemplo
	fakeBeerIdStr   string = "1"   // ID como string para path params
	fakeQuantityStr string = "1"   // Cantidad como string para query params
)

// HTTPContext encapsula los parámetros para simular peticiones HTTP en Echo.
type HTTPContext struct {
	Req         *http.Request              // Petición simulada
	Res         *httptest.ResponseRecorder // Grabador de respuesta
	EchoContext echo.Context               // Contexto Echo construido
}

// SetupHTTPContext prepara un contexto Echo con:
// - Método, ruta y cuerpo (si body != nil).
// - Parámetros de ruta (pathParams).
// - Query params (queryParams).
// - Header Content-Type (mediaType).
func SetupHTTPContext(
	method, routePattern string,
	body interface{},
	pathParams map[string]string,
	queryParams map[string]string,
	mediaType string,
) HTTPContext {
	// 1) Serializar body si existe
	var reader io.Reader
	if body != nil {
		payload, _ := json.Marshal(body)
		reader = bytes.NewReader(payload)
	}

	// 2) Crear servidor Echo, request y response recorder
	e := echo.New()
	req := httptest.NewRequest(method, routePattern, reader)
	rec := httptest.NewRecorder()
	// Establecer Content-Type si se indicó
	if mediaType != "" {
		req.Header.Set(echo.HeaderContentType, mediaType)
	}

	// 3) Agregar query params a la URL
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// 4) Crear contexto Echo
	ctx := e.NewContext(req, rec)

	// 5) Asignar parámetros de ruta de una vez
	if len(pathParams) > 0 {
		keys, vals := make([]string, 0, len(pathParams)), make([]string, 0, len(pathParams))
		for k, v := range pathParams {
			keys = append(keys, k)
			vals = append(vals, v)
		}
		ctx.SetParamNames(keys...)
		ctx.SetParamValues(vals...)
	}

	return HTTPContext{Req: req, Res: rec, EchoContext: ctx}
}

// Test_beerHandler_GetOneHandler prueba todas las rutas de GetOneHandler:
// - Respuesta exitosa (200) con JSON correcto.
// - BadRequest (400) cuando beerID es inválido o faltante.
// - NotFound (404) cuando el servicio devuelve ErrBeerNotFound.
// - Internal Server Error (500) cuando ocurre otro error.
func Test_beerHandler_GetOneHandler(t *testing.T) {
	mockService := _mocksService.NewMockBeerService(t)
	ctx := context.Background()
	handler := NewBeerHandler(mockService)

	t.Run("Successful Response", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			nil,
			"")

		expectedDataResponse := model.BeersResponse{
			ID:       fakeBeerIdUint,
			Name:     "Gulden Draak",
			Brewery:  "Blót",
			Country:  "BE",
			Price:    6.50,
			Currency: "EUR",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		mockService.On("GetBeerById", ctx, fakeBeerIdUint).Return(expectedDataResponse, nil).Once()

		res := httpContext.Res
		err := handler.GetOneHandler(httpContext.EchoContext)

		expectedResponse := `{
			  "id": 1,
			  "name": "Gulden Draak",
			  "brewery": "Blót",
			  "country": "BE",
			  "price": 6.50,
			  "currency": "EUR",
			  "created_at": "2023-01-01T00:00:00Z",
			  "updated_at": "2023-01-01T00:00:00Z"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Bad Request", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"": ""},
			map[string]string{},
			"",
		)

		res := httpContext.Res
		err := handler.GetOneHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "Invalid Beer ID"	
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
	})

	t.Run("Successful Not Found", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			map[string]string{},
			"",
		)

		mockService.On("GetBeerById", ctx, fakeBeerIdUint).Return(model.BeersResponse{}, repository.ErrBeerNotFound).Once()

		res := httpContext.Res
		err := handler.GetOneHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "beer not found"
		}
`
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Internal Server Error", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			map[string]string{},
			"",
		)

		mockService.On("GetBeerById", ctx, fakeBeerIdUint).Return(model.BeersResponse{}, assert.AnError).Once()

		res := httpContext.Res
		err := handler.GetOneHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "assert.AnError general error for testing"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

}

// Test_beerHandler_GetOneBoxPriceHandler cubre los casos de uso de GetOneBoxPriceHandler:
// - Respuesta exitosa (200) con PriceResponse.
// - BadRequest (400) para parámetros faltantes o inválidos.
// - Internal Server Error (500) ante error del servicio.
func Test_beerHandler_GetOneBoxPriceHandler(t *testing.T) {
	mockService := _mocksService.NewMockBeerService(t)
	ctx := context.Background()
	handler := NewBeerHandler(mockService)

	t.Run("Successful Response", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			map[string]string{
				"currency": fakeCurrencyStr,
				"quantity": fakeQuantityStr,
			},
			"",
		)

		expectedDataResponse := model.PriceResponse{
			PriceTotal:  33.68235,
			CurrencyPay: fakeCurrencyStr,
		}

		mockService.On("GetOneBoxPrice", ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt).
			Return(expectedDataResponse, nil).Once()

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"price_total": 33.68235,
			"currency_pay": "USD"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Bad Request for BeerID", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"": ""},
			nil,
			"",
		)

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "beerID is required"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
	})

	t.Run("Successful Bad Request for BeerID Not Numeric", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": "not a number"},
			map[string]string{
				"currency": fakeCurrencyStr,
				"quantity": fakeQuantityStr,
			},
			"",
		)

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "beerID must be numeric"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
	})

	t.Run("Successful Bad Request for currency", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			nil,
			"",
		)

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "currency is required and must be at least 3 characters"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
	})

	t.Run("Successful Bad Request for Quantity", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			map[string]string{
				"currency": fakeCurrencyStr,
				"quantity": "not a number",
			},
			"",
		)

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "quantity must be a positive integer"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

	})

	t.Run("Successful Not Found", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			map[string]string{
				"currency": fakeCurrencyStr,
				"quantity": fakeQuantityStr,
			},
			"",
		)

		expectedDataResponse := model.PriceResponse{
			PriceTotal:  33.68235,
			CurrencyPay: fakeCurrencyStr,
		}

		mockService.On("GetOneBoxPrice", ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt).
			Return(expectedDataResponse, nil).Once()

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"price_total": 33.68235,
			"currency_pay": "USD"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Internal Server Error", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/:beerID",
			nil,
			map[string]string{"beerID": fakeBeerIdStr},
			map[string]string{
				"currency": fakeCurrencyStr,
				"quantity": fakeQuantityStr,
			},
			"",
		)

		mockService.On("GetOneBoxPrice", ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt).
			Return(model.PriceResponse{}, assert.AnError).Once()

		res := httpContext.Res
		err := handler.GetOneBoxPriceHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "assert.AnError general error for testing"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})
}

// Test_beerHandler_CreateHandler valida CreateHandler:
// - 201 Created cuando la creación es exitosa.
// - 400 BadRequest para JSON mal formado.
// - 422 Unprocessable Entity si Validation falla.
// - 409 Conflict si el servicio indica entidad existente.
// - 500 Internal Server Error para otros errores.
func Test_beerHandler_CreateHandler(t *testing.T) {
	mockService := _mocksService.NewMockBeerService(t)
	ctx := context.Background()
	handler := NewBeerHandler(mockService)
	bodyRequest := model.BeersRequest{
		ID:       fakeBeerIdUint,
		Name:     "Gulden Draak",
		Brewery:  "Blót",
		Country:  "BE",
		Price:    6.50,
		Currency: "EUR",
	}

	t.Run("Successful Response", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodPost,
			"/:beerID",
			bodyRequest,
			nil,
			nil,
			echo.MIMEApplicationJSON,
		)

		mockService.On("CreateBeerWithId", ctx, &bodyRequest).Return(nil).Once()

		res := httpContext.Res
		err := handler.CreateHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "Beer created"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Bad Request", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodPost,
			"/:beerID",
			bodyRequest,
			nil,
			nil,
			"",
		)

		res := httpContext.Res
		err := handler.CreateHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "Invalid JSON format or poorly formatted fields"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
	})

	t.Run("Successful Bad Request Validate Fields", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodPost,
			"/:beerID",
			nil,
			nil,
			nil,
			echo.MIMEApplicationJSON,
		)

		res := httpContext.Res
		err := handler.CreateHandler(httpContext.EchoContext)

		expectedResponse := `{
			"brewery_required": "brewery is required",
			"country_required": "country is required",
			"currency_required": "currency is required and it has to be a valid currency code",
			"id_required": "Id is required or id invalid",
			"name_required": "name is required",
			"price_required": "price is required and must be greater than zero"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
	})

	t.Run("Successful Conflict Beer Already Exists", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodPost,
			"/:beerID",
			bodyRequest,
			map[string]string{},
			map[string]string{},
			echo.MIMEApplicationJSON,
		)

		mockService.On("CreateBeerWithId", ctx, &bodyRequest).Return(service.ErrBeerAlreadyExists).Once()

		res := httpContext.Res
		err := handler.CreateHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "beer ID already exists"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Internal Server Error", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodPost,
			"/:beerID",
			bodyRequest,
			nil,
			nil,
			echo.MIMEApplicationJSON,
		)

		mockService.On("CreateBeerWithId", ctx, &bodyRequest).Return(assert.AnError).Once()

		res := httpContext.Res
		err := handler.CreateHandler(httpContext.EchoContext)

		expectedResponse := `{
			"message": "assert.AnError general error for testing"
		}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())

		// Verify that the mocks were called as expected
		mockService.AssertExpectations(t)
	})

}
