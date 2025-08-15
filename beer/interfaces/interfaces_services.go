package interfaces

import (
	"context"

	"go_trainning/beer-api/beer/model"
)

type BeerService interface {
    GetAllBeers(ctx context.Context) ([]model.Beers, error)
    GetBeerById(ctx context.Context, id uint) (model.Beers, error)
    // CreateBeerWithId crea una nueva cerveza basada en los datos de BeersRequest
    // Antes de persistir, asigna las marcas de tiempo CreatedAt y UpdatedAt.
    //
    // Parámetros:
	// - ctx: contexto para control de tiempo de espera y cancelación
	// - beersReq: datos de la cerveza a crear, provenientes de la petición HTTP
    //
    // Retorna:
	// - error: en caso de fallo en la creación
	CreateBeerWithId(ctx context.Context, beers *model.BeersRequest) error
    GetOneBoxPrice(ctx context.Context, id uint, currencyPay string, quantity int) (model.PriceResponse, error)
}