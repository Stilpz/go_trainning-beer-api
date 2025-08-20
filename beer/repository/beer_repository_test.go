// Package repository_test contiene pruebas unitarias para la implementación
// de BeerRepository usando SQLMock y testify.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	// "github.com/samuskitchen/beer-api/beer/model"
	"go_trainning/beer-api/beer/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	fakeBeerIdUint uint = 1 // ID de cerveza de ejemplo para pruebas
)

// dataBeers retorna una muestra de datos de cervezas para usar en las pruebas.
func dataBeers() []model.Beers {
	now := time.Now()

	return []model.Beers{
		{
			ID:        fakeBeerIdUint,
			Name:      "Gulden Draak",
			Brewery:   "Blót",
			Country:   "BE",
			Price:     6.50,
			Currency:  "EUR",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uint(2),
			Name:      "Club Colombia",
			Brewery:   "Bavaria",
			Country:   "CO",
			Price:     3.483,
			Currency:  "COP",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// Test_beerRepository_GetBeerById verifica los diferentes flujos de GetBeerById:
// - Éxito al retornar una fila válida.
// - Error de consulta SQL.
// - Registro no encontrado (no rows).
func Test_beerRepository_GetBeerById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBeerRepository(db)
	ctx := context.Background()
	beerTest := dataBeers()[0]

	t.Run("Success SQL", func(tt *testing.T) {
		// Preparar filas de respuesta
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at", "updated_at"}).
			AddRow(
				beerTest.ID,
				beerTest.Name,
				beerTest.Brewery,
				beerTest.Country,
				beerTest.Price,
				beerTest.Currency,
				beerTest.CreatedAt,
				beerTest.UpdatedAt,
			)

		// Esperar llamada a QueryContext con el SQL correcto
		mock.ExpectQuery(regexp.QuoteMeta(selectBeerById)).
			WithArgs(fakeBeerIdUint).
			WillReturnRows(rows)

		gotBeer, errRepo := repo.GetBeerById(ctx, fakeBeerIdUint)

		assert.NoError(t, errRepo)
		assert.Equal(t, beerTest, gotBeer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error SQL", func(tt *testing.T) {
		// Simular error en la ejecución de la query
		mock.ExpectQuery(regexp.QuoteMeta(selectBeerById)).WithArgs(fakeBeerIdUint).WillReturnError(assert.AnError)
		gotBeer, errRepo := repo.GetBeerById(ctx, fakeBeerIdUint)

		assert.Error(t, errRepo)
		assert.Empty(t, gotBeer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found SQL", func(tt *testing.T) {
		// Simular RowError de sql.ErrNoRows
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at", "updated_at"}).
			AddRow(
				beerTest.ID,
				beerTest.Name,
				beerTest.Brewery,
				beerTest.Country,
				beerTest.Price,
				beerTest.Currency,
				beerTest.CreatedAt,
				beerTest.UpdatedAt,
			).RowError(0, sql.ErrNoRows)

		mock.ExpectQuery(regexp.QuoteMeta(selectBeerById)).WithArgs(fakeBeerIdUint).WillReturnRows(rows)
		gotBeer, errRepo := repo.GetBeerById(ctx, fakeBeerIdUint)

		// Debe mapear sql.ErrNoRows a ErrBeerNotFound
		assert.Error(t, errRepo)
		assert.Equal(t, errRepo, ErrBeerNotFound)
		assert.Empty(tt, model.Beers{}, gotBeer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_beerRepository_GetAllBeers valida la funcionalidad de obtener todas las cervezas.
// TODO: implementar casos de prueba para GetAllBeers.
func Test_beerRepository_GetAllBeers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBeerRepository(db)
	ctx := context.Background()
	beersTest := dataBeers()

	t.Run("Success SQL", func(tt *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at", "updated_at"})
		for _, beer := range beersTest {
			rows.AddRow(beer.ID, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency, beer.CreatedAt, beer.UpdatedAt)
		}
		mock.ExpectQuery(regexp.QuoteMeta(selectAllBeers)).WillReturnRows(rows)

		gotBeers, errRepo := repo.GetAllBeers(ctx)
		assert.NoError(t, errRepo)
		assert.Equal(t, beersTest, gotBeers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error SQL", func(tt *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectAllBeers)).WillReturnError(assert.AnError)
		gotBeers, errRepo := repo.GetAllBeers(ctx)
		assert.Error(t, errRepo)
		assert.Empty(t, gotBeers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No Results", func(tt *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at", "updated_at"})
		mock.ExpectQuery(regexp.QuoteMeta(selectAllBeers)).WillReturnRows(rows)
		gotBeers, errRepo := repo.GetAllBeers(ctx)
		assert.NoError(t, errRepo)
		assert.Empty(t, gotBeers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_beerRepository_CreateBeerWithId verifica los flujos de inserción:
// - Éxito al preparar statement e insertar fila.
// - Error al preparar statement.
// - Error al escanear el row result.
func Test_beerRepository_CreateBeerWithId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBeerRepository(db)
	ctx := context.Background()
	beerTest := dataBeers()[0]

	t.Run("Success Prepare Statement", func(tt *testing.T) {
		// Mock de PrepareContext y QueryRowContext
		prep := mock.ExpectPrepare(regexp.QuoteMeta(insertBeerWithId))
		prep.ExpectQuery().
			WithArgs(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt, beerTest.UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(beerTest.ID))

		errRepo := repo.CreateBeerWithId(ctx, &beerTest)

		assert.NoError(t, errRepo)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error Prepare Statement", func(tt *testing.T) {
		// Simular error al preparar el statement
		wantErr := errors.New("prepare error")
		mock.
			ExpectPrepare(regexp.QuoteMeta(insertBeerWithId)).
			WillReturnError(wantErr)

		errRepo := repo.CreateBeerWithId(ctx, &beerTest)
		assert.Error(tt, errRepo)
		assert.Equal(tt, wantErr, errRepo)

		assert.NoError(tt, mock.ExpectationsWereMet())
	})

	t.Run("Error Scan Row", func(tt *testing.T) {
		// Simular RowError durante Scan
		wantErr := errors.New("iteration error")

		prep := mock.ExpectPrepare(regexp.QuoteMeta(insertBeerWithId))
		prep.
			ExpectQuery().
			WithArgs(
				beerTest.ID, beerTest.Name, beerTest.Brewery,
				beerTest.Country, beerTest.Price, beerTest.Currency,
				beerTest.CreatedAt, beerTest.UpdatedAt,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(beerTest.ID).
					RowError(0, wantErr),
			)

		errRepo := repo.CreateBeerWithId(ctx, &beerTest)

		assert.Error(t, errRepo)
		assert.Equal(t, wantErr, errRepo)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
