package interfaces

import (
	"context"
	"go_trainning/beer-api/beer/model"
)

type BeerRepository interface {
	GetAllBeers(ctx context.Context) ([]model.Beers, error)
	GetBeerById(ctx context.Context, id uint) (model.Beers, error)
	CreateBeerWithId(ctx context.Context, id uint, beer model.Beers) error
}