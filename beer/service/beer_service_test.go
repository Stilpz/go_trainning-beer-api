// Package service_test contiene pruebas unitarias para la capa de servicio BeerService.
// Utiliza testify para aserciones y mocks de las interfaces de repositorio y cliente externo.
package service

import (
	"context"
	"testing"
	"time"

	externalModel "go_trainning/beer-api/beer/external/model"
	_mockInterfaces "go_trainning/beer-api/beer/mocks/interfaces"
	"go_trainning/beer-api/beer/model"
	"go_trainning/beer-api/beer/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	fakeBeerIdUint  uint   = 1     // ID de cerveza de ejemplo
	fakeQuantityInt int    = 6     // Cantidad de ejemplo para GetOneBoxPrice
	fakeCurrencyStr string = "USD" // Moneda de ejemplo
)

// dataBeers prepara datos de modelo.Beers para usar en las pruebas.
func dataBeers() []model.Beers {
	now := time.Now()
	return []model.Beers{
		{ID: fakeBeerIdUint, Name: "Gulden Draak", Brewery: "Blót", Country: "BE", Price: 6.50, Currency: "EUR", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Club Colombia", Brewery: "Bavaria", Country: "CO", Price: 3.483, Currency: "COP", CreatedAt: now, UpdatedAt: now},
	}
}

// dataBeersRequest construye un modelo.BeersRequest para pruebas de creación.
func dataBeersRequest() model.BeersRequest {
	return model.BeersRequest{ID: fakeBeerIdUint, Name: "Gulden Draak", Brewery: "Blót", Country: "BE", Price: 6.50, Currency: "EUR"}
}

// dataBeersResponse retorna la representación BeersResponse correspondiente.
func dataBeersResponse() []model.BeersResponse {
	return []model.BeersResponse{
		{ID: fakeBeerIdUint, Name: "Gulden Draak", Brewery: "Blót", Country: "BE", Price: 6.5, Currency: "EUR"},
		{ID: 2, Name: "Club Colombia", Brewery: "Bavaria", Country: "CO", Price: 3.483, Currency: "COP"},
	}
}

// dataCurrencyResponse prepara un externalModel.CurrencyConversionResponse simulado.
func dataCurrencyResponse() externalModel.CurrencyConversionResponse {
       return externalModel.CurrencyConversionResponse{
	       Success: true,
	       Query: externalModel.CurrencyQuery{
		       From:   "USD",
		       To:     "EUR",
		       Amount: 39, // 6.5 * 6
	       },
	       Info: externalModel.InfoResponse{Timestamp: 1234567890, Rate: 1.15553},
	       Result:  45.07, // 39 * 1.15553 redondeado a dos decimales
       }
}

// dataPriceResponse retorna el modelo.PriceResponse esperado tras la conversión.
func dataPriceResponse() model.PriceResponse {
	return model.PriceResponse{CurrencyPay: fakeCurrencyStr, PriceTotal: 45.07}
}

// Test_beerService_GetAllBeers placeholder para pruebas de GetAllBeers.
// TODO: implementar casos de prueba.
func Test_beerService_GetAllBeers(t *testing.T) {
    ctx := context.Background()

    t.Run("success", func(t *testing.T) {
        beers := dataBeers()
        beersResponse := make([]model.BeersResponse, 0, len(beers))
        for _, b := range beers {
            beersResponse = append(beersResponse, b.ToBeersResponse())
        }

        mockRepository := _mockInterfaces.NewMockBeerRepository(t)
        mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
        service := NewBeerService(mockRepository, mockExternal)

        mockRepository.On("GetAllBeers", ctx).Return(beers, nil)

        gotBeers, errService := service.GetAllBeers(ctx)
        assert.NoError(t, errService)
        assert.Equal(t, beersResponse, gotBeers)
        mockRepository.AssertExpectations(t)
    })

    t.Run("error repository", func(t *testing.T) {
        mockRepository := _mockInterfaces.NewMockBeerRepository(t)
        mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
        service := NewBeerService(mockRepository, mockExternal)

        mockRepository.On("GetAllBeers", ctx).Return(nil, assert.AnError)

        gotBeers, errService := service.GetAllBeers(ctx)
        assert.Error(t, errService)
        assert.Nil(t, gotBeers)
        mockRepository.AssertExpectations(t)
    })
}
// Test_beerService_GetBeerById verifica los flujos de GetBeerById:
// - Éxito devolviendo BeersResponse
// - Error del repositorio
func Test_beerService_GetBeerById(t *testing.T) {
		ctx := context.Background()

		t.Run("success", func(t *testing.T) {
			now := time.Now()
			beerTest := model.Beers{
				ID:        fakeBeerIdUint,
				Name:      "Gulden Draak",
				Brewery:   "Blót",
				Country:   "BE",
				Price:     6.50,
				Currency:  "EUR",
				CreatedAt: now,
				UpdatedAt: now,
			}
			beerResponseTest := model.BeersResponse{
				ID:        fakeBeerIdUint,
				Name:      "Gulden Draak",
				Brewery:   "Blót",
				Country:   "BE",
				Price:     6.5,
				Currency:  "EUR",
				CreatedAt: now,
				UpdatedAt: now,
			}

			mockRepository := _mockInterfaces.NewMockBeerRepository(t)
			mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
			service := NewBeerService(mockRepository, mockExternal)

			mockRepository.On("GetBeerById", ctx, fakeBeerIdUint).Return(beerTest, nil)

			gotBeer, errService := service.GetBeerById(ctx, fakeBeerIdUint)
			assert.NoError(t, errService)
			assert.NotEmpty(t, gotBeer)
			assert.Equal(t, beerResponseTest.ID, gotBeer.ID)
			assert.Equal(t, beerResponseTest.Name, gotBeer.Name)
			assert.Equal(t, beerResponseTest.Brewery, gotBeer.Brewery)
			assert.Equal(t, beerResponseTest.Country, gotBeer.Country)
			assert.Equal(t, beerResponseTest.Price, gotBeer.Price)
			assert.Equal(t, beerResponseTest.Currency, gotBeer.Currency)
			assert.WithinDuration(t, beerResponseTest.CreatedAt, gotBeer.CreatedAt, time.Second)
			assert.WithinDuration(t, beerResponseTest.UpdatedAt, gotBeer.UpdatedAt, time.Second)

			mockRepository.AssertExpectations(t)
		})

	t.Run("error repository", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		wantErr := repository.ErrBeerNotFound
		mockRepository.On("GetBeerById", ctx, fakeBeerIdUint).Return(model.Beers{}, wantErr)

		gotBeer, errService := service.GetBeerById(ctx, fakeBeerIdUint)

		assert.Error(t, errService)
		assert.Empty(t, gotBeer)
		assert.Equal(t, wantErr, errService)

		mockRepository.AssertExpectations(t)
	})
}

// Test_beerService_GetOneBoxPrice verifica los flujos de GetOneBoxPrice:
// - Éxito
// - Error al obtener la cerveza
// - Error en repositorio general
// - Error en conversión de moneda
func Test_beerService_GetOneBoxPrice(t *testing.T) {
	ctx := context.Background()
	beerTest := dataBeers()[0]

	t.Run("success", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		mockRepository.On("GetBeerById", ctx, fakeBeerIdUint).Return(beerTest, nil)
		mockExternal.On("GetExchangeCurrency", fakeCurrencyStr, beerTest.Currency, beerTest.Price*float64(fakeQuantityInt)).
			Return(dataCurrencyResponse(), nil)

		gotPrice, errService := service.GetOneBoxPrice(ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt)

		assert.NoError(t, errService)
		assert.NotEmpty(t, gotPrice)
		assert.Equal(t, dataPriceResponse(), gotPrice)

		mockRepository.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
	})

	t.Run("error beer not found", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		wantErr := repository.ErrBeerNotFound
		mockRepository.On("GetBeerById", ctx, fakeBeerIdUint).Return(model.Beers{}, wantErr)

		gotPrice, errService := service.GetOneBoxPrice(ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt)

		assert.Error(t, errService)
		assert.Empty(t, gotPrice)
		assert.Equal(t, wantErr, errService)

		mockRepository.AssertExpectations(t)
	})

	t.Run("error repository", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		wantErr := assert.AnError
		mockRepository.On("GetBeerById", ctx, fakeBeerIdUint).Return(model.Beers{}, wantErr)

		gotPrice, errService := service.GetOneBoxPrice(ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt)

		assert.Error(t, errService)
		assert.Empty(t, gotPrice)
		assert.Equal(t, wantErr, errService)

		mockRepository.AssertExpectations(t)
	})

	t.Run("success exchange currency", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		wantErr := assert.AnError
		mockRepository.On("GetBeerById", ctx, fakeBeerIdUint).Return(beerTest, nil)
		mockExternal.On("GetExchangeCurrency", fakeCurrencyStr, beerTest.Currency, beerTest.Price*float64(fakeQuantityInt)).
			Return(externalModel.CurrencyConversionResponse{}, wantErr)

		gotPrice, errService := service.GetOneBoxPrice(ctx, fakeBeerIdUint, fakeCurrencyStr, fakeQuantityInt)

		assert.Error(t, errService)
		assert.Empty(t, gotPrice)
		assert.Equal(t, wantErr, errService)

		mockRepository.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
	})
}

// Test_beerService_CreateBeerWithId cubre los flujos de creación:
// - Éxito al no encontrar cerveza previa y crear
// - Error si cerveza ya existe
// - Error al consultar repositorio
// - Error al insertar
func Test_beerService_CreateBeerWithId(t *testing.T) {
	ctx := context.Background()
	beerTest := dataBeers()[0]
	beerRequest := dataBeersRequest()

	t.Run("success", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		mockRepository.
			On("GetBeerById", ctx, fakeBeerIdUint).
			Return(model.Beers{}, repository.ErrBeerNotFound)

		// CreateBeerWithId: comprobamos solo los campos esenciales
		mockRepository.
			On("CreateBeerWithId", ctx,
				mock.MatchedBy(func(arg *model.Beers) bool {
					return arg.ID == beerRequest.ID &&
						arg.Name == beerRequest.Name &&
						arg.Brewery == beerRequest.Brewery &&
						arg.Country == beerRequest.Country &&
						arg.Price == beerRequest.Price &&
						arg.Currency == beerRequest.Currency &&
						// los timestamps se asignaron con time.Now(),
						// así que sólo verificamos que sean recientes
						time.Since(arg.CreatedAt) < time.Second &&
						time.Since(arg.UpdatedAt) < time.Second
				}),
			).
			Return(nil)

		errService := service.CreateBeerWithId(ctx, &beerRequest)

		assert.NoError(t, errService)
		mockRepository.AssertExpectations(t)
	})

	t.Run("error beer already exists", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		mockRepository.
			On("GetBeerById", ctx, fakeBeerIdUint).
			Return(beerTest, nil)

		errService := service.CreateBeerWithId(ctx, &beerRequest)

		assert.Error(t, errService)
		mockRepository.AssertExpectations(t)
	})

	t.Run("error repository get beer", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		wantErr := assert.AnError
		mockRepository.
			On("GetBeerById", ctx, fakeBeerIdUint).
			Return(model.Beers{}, wantErr)

		errService := service.CreateBeerWithId(ctx, &beerRequest)

		assert.Error(t, errService)
		assert.Equal(t, wantErr, errService)
		mockRepository.AssertExpectations(t)
	})

	t.Run("error repository create beer", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockExternal := _mockInterfaces.NewMockCurrencyLayer(t)
		service := NewBeerService(mockRepository, mockExternal)

		wantErr := assert.AnError
		mockRepository.
			On("GetBeerById", ctx, fakeBeerIdUint).
			Return(model.Beers{}, repository.ErrBeerNotFound)

		mockRepository.
			On("CreateBeerWithId", ctx,
				mock.MatchedBy(func(arg *model.Beers) bool {
					return arg.ID == beerRequest.ID &&
						arg.Name == beerRequest.Name &&
						arg.Brewery == beerRequest.Brewery &&
						arg.Country == beerRequest.Country &&
						arg.Price == beerRequest.Price &&
						arg.Currency == beerRequest.Currency &&
						// los timestamps se asignaron con time.Now(),
						// así que sólo verificamos que sean recientes
						time.Since(arg.CreatedAt) < time.Second &&
						time.Since(arg.UpdatedAt) < time.Second
				}),
			).
			Return(wantErr)

		errService := service.CreateBeerWithId(ctx, &beerRequest)

		assert.Error(t, errService)
		assert.Equal(t, wantErr, errService)
		mockRepository.AssertExpectations(t)

	})
}
