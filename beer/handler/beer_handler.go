// Package handler define los controladores HTTP (handlers) para la entidad Beer.
// Se encargan de recibir las solicitudes, validarlas y delegar la lógica
// de negocio al servicio correspondiente. Además de manejar respuestas y errores.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go_trainning/beer-api/beer/interfaces"
	"go_trainning/beer-api/beer/model"
	"go_trainning/beer-api/beer/repository"
	"go_trainning/beer-api/beer/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// errorResponse estructura para las respuestas genericas de error
type errorResponse struct {
	Message string `json:"message"`
}

// beerHandler implementa BeerHandler y encapsula el servicio de Beer.
type beerHandler struct {
	beerService interfaces.BeerService
}

// BeerHandler agrupa los métodos handler para endpoints de Beer.
type BeerHandler interface {
	GetAllBeersHandler(c echo.Context) error    // Listar todas las cervezas
	GetOneHandler(c echo.Context) error         // Obtener una cerveza por ID
	CreateHandler(c echo.Context) error         // Crear una nueva cerveza
	GetOneBoxPriceHandler(c echo.Context) error // Calcular precio de una caja de cervezas
}

// NewBeerHandler construye un BeerHandler con la implementación del servicio.
func NewBeerHandler(service interfaces.BeerService) BeerHandler {
	return &beerHandler{beerService: service}
}

// Greeting devuelve un saludo simple para comprobar el estado del servicio.
func Greeting(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

// GetAllBeersHandler responde con la lista de todas las cervezas.
// Actualmente devuelve un placeholder.
// func (bh *beerHandler) GetAllBeersHandler(c echo.Context) error {
// 	// TODO: implementar llamada a bh.beerService.GetAllBeers
// 	return c.JSON(http.StatusOK, "all beers")
// }

func (bh *beerHandler) GetAllBeersHandler(c echo.Context) error {
	// TODO: implementar llamada a bh.beerService.GetAllBeers
	ctx := c.Request().Context()
	beers, err := bh.beerService.GetAllBeers(ctx)
	if err != nil {
		log.Error().Msgf("error GetAllBeers: %v", err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, beers)
}

// GetOneHandler obtiene una cerveza por su ID, validando el parámetro y devolviendo JSON.
// @Description Obtiene una cerveza por su ID
// @Tags Beer
// @ID GetOneHandler
// @Param beerID path string true "Beer ID"
// @Success 200 {object} model.BeersResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /{beerID} [GET]
func (bh *beerHandler) GetOneHandler(c echo.Context) error {
	ctx := c.Request().Context()

	// Leer y validar parámetro path beerID
	idParam := c.Param("beerID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error().Msgf("invalid beer ID: %v", err)
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "Invalid Beer ID"})
	}

	// Llamar al servicio
	response, err := bh.beerService.GetBeerById(ctx, uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrBeerNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{Message: repository.ErrBeerNotFound.Error()})
		}

		log.Error().Msgf("error getting beer by ID: %v", err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
	}

	// Devolver resultado
	return c.JSON(http.StatusOK, response)
}

// CreateHandler recibe un JSON con datos de cerveza, lo valida y crea la entidad.
// @Description Recibe un JSON con datos de cerveza
// @Tags Beer
// @Accept application/json
// @ID CreateHandler
// @Param beer-body body model.BeersRequest true "request body"
// @Success 201 {string} Beer created
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 422 {object} map[string]string
// @Failure 500 {object} errorResponse
// @Router / [POST]
func (bh *beerHandler) CreateHandler(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.BeersRequest

	// Bindear JSON a struct
	if err := c.Bind(&req); err != nil {
		log.Error().Msgf("invalid JSON format: %v", err)
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "Invalid JSON format or poorly formatted fields"})
	}

	// Validar campos
	if errs := req.Validate(); len(errs) > 0 {
		log.Error().Msgf("validation failed for BeersRequest: %v", errs)
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	// Crear cerveza
	if err := bh.beerService.CreateBeerWithId(ctx, &req); err != nil {
		if errors.Is(err, service.ErrBeerAlreadyExists) {
			return c.JSON(http.StatusConflict, errorResponse{Message: service.ErrBeerAlreadyExists.Error()})
		}
		
		log.Error().Msgf("error creating beer: %v", err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
	}

	// Responder con código 201
	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "Beer created"})
}

// GetOneBoxPriceHandler calcula y devuelve el precio total de una caja de cervezas.
// - Parámetros path: beerID
// - Query params: currencyPay (código ISO), quantity (opcional, default 6)
// @Description Calcula y devuelve el precio total de una caja de cervezas
// @Tags Currency
// @ID GetOneBoxPriceHandler
// @Param beerID path string true "Beer ID"
// @Param currency query string true "Currency Pay"
// @Param quantity query int true "Quantity"
// @Success 200 {object} model.PriceResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /{beerID}/box-price [GET]
func (bh *beerHandler) GetOneBoxPriceHandler(c echo.Context) error {
	ctx := c.Request().Context()
	// Valores por defecto
	quantity := 6

	// Leer parámetros
	idStr := c.Param("beerID")
	currencyPay := c.QueryParam("currency")
	qtyStr := c.QueryParam("quantity")

	// Validaciones básicas
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "beerID is required"})
	}
	if len(currencyPay) < 3 {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "currency is required and must be at least 3 characters"})
	}
	if qtyStr != "" {
		if qtyParam, err := strconv.Atoi(qtyStr); err != nil || qtyParam <= 0 {
			return c.JSON(http.StatusBadRequest, errorResponse{Message: "quantity must be a positive integer"})
		} else {
			quantity = qtyParam
		}
	}

	// Parsear ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "beerID must be numeric"})
	}

	// Llamar al servicio de negocio
	priceResp, err := bh.beerService.GetOneBoxPrice(ctx, uint(id), currencyPay, quantity)
	if err != nil {
		if errors.Is(err, repository.ErrBeerNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{Message: repository.ErrBeerNotFound.Error()})
		}
		
		log.Error().Msgf("error get price: %v", err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
	}

	// Devolver precio calculado
	return c.JSON(http.StatusOK, priceResp)
}