// Package service contiene la implementación de los casos de uso de la entidad Beer.
// Orquesta la interacción entre el repositorio de cervezas y el cliente de conversión
// de moneda para ofrecer funcionalidades de alto nivel.
package service

import (
	"context"
	"errors"
	"math"
	"time"

	"go_trainning/beer-api/beer/interfaces"
	"go_trainning/beer-api/beer/model"
	"go_trainning/beer-api/beer/repository"

	"github.com/rs/zerolog/log"
)

// ErrBeerAlreadyExists es el error sentinela que devuelvo al encontrar la cerveza.
var ErrBeerAlreadyExists = errors.New("beer ID already exists")

// beerService implementa interfaces.BeerService.
// Combina un repositorio de cervezas y un cliente de conversión de moneda.
type beerService struct {
	beerRepository      interfaces.BeerRepository
	clientCurrencyLayer interfaces.CurrencyLayer
}

// NewBeerService crea una nueva instancia de BeerService.
// Parámetros:
//   - beerRepository: implementación de interfaces.BeerRepository.
//   - clientCurrency: implementación de interfaces.CurrencyLayer.
//
// Retorna:
//   - interfaces.BeerService: servicio listo para usarse.
func NewBeerService(
	beerRepository interfaces.BeerRepository,
	clientCurrency interfaces.CurrencyLayer,
) interfaces.BeerService {
	return &beerService{
		beerRepository:      beerRepository,
		clientCurrencyLayer: clientCurrency,
	}
}

func (b *beerService) GetAllBeers(ctx context.Context) ([]model.BeersResponse, error) {
	subLogger := log.With().Str("Method", "BeerService.GetAllBeers").Logger()
	subLogger.Info().Msg("INIT")
	beers, err := b.beerRepository.GetAllBeers(ctx)
	if err != nil {
		subLogger.Error().Msgf("error GetAllBeers repo: %v", err)
		return nil, err
	}
	resp := make([]model.BeersResponse, 0, len(beers))
	for _, v := range beers {
		resp = append(resp, v.ToBeersResponse())
	}
	subLogger.Info().Msg("END_OK")
	return resp, nil
}

// GetBeerById obtiene una cerveza por su ID y devuelve su respuesta formateada.
// Parámetros:
//   - ctx: contexto para control de tiempo de espera y cancelación.
//   - ID: identificador único de la cerveza.
//
// Retorna:
//   - model.BeersResponse: datos de la cerveza formateados para la respuesta.
//   - error: si no existe o falla la consulta.
func (b *beerService) GetBeerById(ctx context.Context, ID uint) (model.BeersResponse, error) {
    subLogger := log.With().Str("Method", "BeerService.GetBeerById").Logger()
    subLogger.Info().Msg("INIT")
    subLogger.Info().Msgf("argument[s] beer_id=%v", ID)

    beer, err := b.beerRepository.GetBeerById(ctx, ID)
    if err != nil {
        return model.BeersResponse{}, err
    }

    subLogger.Info().Msgf("END_OK | beer_id=%v", ID)
    return beer.ToBeersResponse(), nil
}

// CreateBeerWithId crea una nueva cerveza basada en los datos de BeersRequest.
func (b *beerService) CreateBeerWithId(ctx context.Context, beersReq *model.BeersRequest) error {
	subLogger := log.With().Str("Method", "BeerService.CreateBeerWithId").Logger()
	subLogger.Info().Msg("INIT")
	subLogger.Info().Msgf("argument[s] beer_id=%v", beersReq.ID)

	// Validar existencia previa
	existing, err := b.beerRepository.GetBeerById(ctx, beersReq.ID)
	if err != nil && !errors.Is(err, repository.ErrBeerNotFound) {
		subLogger.Error().Msgf("error fetching beer: %v", err)
		return err
	}

	if existing != (model.Beers{}) {
		return ErrBeerAlreadyExists
	}
	
	// Generar timestamps actuales
	now := time.Now()

	// Mapear BeersRequest a entidad persistible model.Beers
	beerEntity := model.Beers{
		ID:        beersReq.ID,
		Name:      beersReq.Name,
		Brewery:   beersReq.Brewery,
		Country:   beersReq.Country,
		Price:     beersReq.Price,
		Currency:  beersReq.Currency,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Delegar creación al repositorio
	if err := b.beerRepository.CreateBeerWithId(ctx, &beerEntity); err != nil {
		subLogger.Error().Msgf("error CreateBeerWithId: %v", err)
		return err
	}

	subLogger.Info().Msgf("END_OK | beer_id=%v", beersReq.ID)
	return nil
}

// GetOneBoxPrice calcula el precio total de una caja de cervezas,
// convirtiendo el precio unitario desde la moneda original a la moneda destino.
// Parámetros:
//   - ctx: contexto para control de tiempo de espera y cancelación.
//   - id: ID de la cerveza a consultar.
//   - currencyPay: código ISO de la moneda destino (p.ej., "USD").
//   - quantity: número de unidades en la caja.
//
// Retorna:
//   - model.PriceResponse: precio total convertido y redondeado.
//   - error: en caso de fallo en cualquiera de los pasos.
func (b *beerService) GetOneBoxPrice(
	ctx context.Context,
	ID uint,
	currencyPay string,
	quantity int,
) (model.PriceResponse, error) {
	subLogger := log.With().Str("Method", "BeerService.GetOneBoxPrice").Logger()
	subLogger.Info().Msg("INIT")
	subLogger.Info().Msgf("argument[s] beer_id=%v, currency_pay=%v, quantity=%v", ID, currencyPay, quantity)

	// 1. Obtener datos de la cerveza
	beer, err := b.beerRepository.GetBeerById(ctx, ID)
	if err != nil {
		subLogger.Error().Msgf("error GetBeerById: %v", err)
		return model.PriceResponse{}, err
	}

	// 2. Calcular el total base en la moneda original
	baseTotal := beer.Price * float64(quantity)

	// 3. Obtener conversión usando el modelo corregido (el método debe usar el modelo adecuado)
	exchange, err := b.clientCurrencyLayer.GetExchangeCurrency(currencyPay, beer.Currency, baseTotal)
	if err != nil {
		subLogger.Error().Msgf("error GetExchangeCurrency: %v", err)
		return model.PriceResponse{}, err
	}

	// 4. Usar el campo Result de la respuesta de conversión
	rounded := math.Round(exchange.Result*100) / 100

	subLogger.Info().Msgf("END_OK | beer_id=%v, currency_pay=%v, quantity=%v", ID, currencyPay, quantity)
	return model.PriceResponse{PriceTotal: rounded, CurrencyPay: currencyPay}, nil
}