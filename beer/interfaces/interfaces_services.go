package interfaces

import (
	"context"
	"go_trainning/beer-api/beer/model"
)

type BeerService interface {
	GetAllBeers(ctx context.Context) ([]model.Beers, error)
	GetBeerById(ctx context.Context, id uint) (model.Beers, error)
	CreateBeerWithId(ctx context.Context, beers *model.Beers) error
	GetOneBoxPrice(ctx context.Context, id uint, currency string, quantity int) (float64, error)
}