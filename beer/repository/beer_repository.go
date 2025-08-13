// Package repository implementa la logica de acceso a datos para las cervezas,
// usando PostgreSQL como almacenamiento. Aqui se define el repositorio concreto
// que satisface la interfaz BeerRepository.
package repository

import (
	"context"
	"database/sql"
	"errors"

	// interfaces define el contrato que debe cumplir el repositorio.
	"go_trainning/beer-api/beer/interfaces"
	// model contiene las estructuras de dominio (e.g., model.Beers).
	"go_trainning/beer-api/beer/model"

	// zerolog para registro estructurado en cada metodo.
	"github.com/rs/zerolog/log"
)

const (
	// selectAllBeers es una consulta que selecciona todas las filas de la tabla cervezas
	selectAllBeers = "SELECT id, \"name\", brewery, country_code, price, currency, created_at, updated_at FROM beers;"

	// selectBeerById es una consulta que selecciona una fila de la tabla de cervezas en función de la identificación dada.
	selectBeerById = "SELECT id, \"name\", brewery, country_code, price, currency, created_at, updated_at FROM beers WHERE id = $1;"

	// insertBeerWithId Es una consulta que inserta una nueva fila en la tabla de cervezas con un id determinado y
	// utiliza los valores datos para id, "nombre", cervecería, país, precio, moneda y fecha de creación
	insertBeerWithId string = "INSERT INTO beers (id, \"name\", brewery, country_code, price, currency, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
)

// ErrBeerNotFound es el error sentinela que devuelvo al no encontrar la cerveza.
var ErrBeerNotFound = errors.New("beer not found")

// beerRepository es la implementacion de BeerRepository que usa
// *sql.DB para comunicarse con PostgreSQL.
type beerRepository struct {
	Conn *sql.DB
}

// NewBeerRepository construye una instancia de BeerRepository usando la conexion dada.
//
// Parametros:
//   - Connection: puntero a sql.DB ya inicializado.
//
// Retorna:
//   - interfaces.BeerRepository: repositorio listo para usarse.
func NewBeerRepository(Connection *sql.DB) interfaces.BeerRepository {
	return &beerRepository{
		Conn: Connection,
	}
}

// GetAllBeers obtiene todas las cervezas registradas en la base de datos.
//
// Parametros:
//   - ctx: contexto para control de tiempo de espera y cancelacion.
//
// Retorna:
//   - []model.Beers: slice con todas las cervezas.
//   - error: en caso de fallo en la consulta o en el escaneo de filas.
func (pb *beerRepository) GetAllBeers(ctx context.Context) ([]model.Beers, error) {
	// Logger con campo Method para rastrear el origen de logs.
	subLogger := log.With().Str("Method", "BeerRepository.GetAllBeers").Logger()
	subLogger.Info().Msg("INIT")

	// Ejecutar la consulta parametrizada definida en selectAllBeers.
	rows, err := pb.Conn.QueryContext(ctx, selectAllBeers)
	if err != nil {
		subLogger.Error().Msgf("error executing query: %v", err)
		return nil, err
	}

	// Asegurar cierre de filas al finalizar.
	defer func ()  {
		if errClose := rows.Close(); errClose != nil {
			subLogger.Error().Msgf("error closing rows: %v", errClose)
		}
	}()

	// Recorrer cada fila y mapear a model.Beers.
	var beers []model.Beers
	for rows.Next() {
		var beerRow model.Beers
		if errScan := rows.Scan(
			&beerRow.ID,
			&beerRow.Name,
			&beerRow.Brewery,
			&beerRow.Country,
			&beerRow.Price,
			&beerRow.Currency,
			&beerRow.CreatedAt,
			&beerRow.UpdatedAt,
		); errScan != nil {
			subLogger.Error().Msgf("error scanning row: %v", errScan)
			return nil, errScan
		}
		beers = append(beers, beerRow)
	}

	subLogger.Info().Msgf("FIN_OK")
	return beers, nil
}

// GetBeerById busca una cerveza por su identificador unico.
//
// Parametros:
//   - ctx: contexto para control de tiempo de espera y cancelacion.
//   - ID: uint que representa la llave primaria de la cerveza.
//
// Retorna:
//   - model.Beers: estructura con los datos de la cerveza.
//   - error: sql.ErrNoRows se traduce a un error "beer not found".
func (pb *beerRepository) GetBeerById(ctx context.Context, ID uint) (model.Beers, error) {
	subLogger := log.With().Str("Method", "BeerRepository.GetBeerById").Logger()
	subLogger.Info().Msg("INIT")
	subLogger.Info().Msgf("argument[s] beer_id=%v", ID)

	// Ejecutar QueryRowContext para un unico resultado.
	row := pb.Conn.QueryRowContext(ctx, selectBeerById, ID)
	var beerScan model.Beers

	// Escanear columnas en la estructura beerScan.
	if errScan := row.Scan(
		&beerScan.ID,
		&beerScan.Name,
		&beerScan.Brewery,
		&beerScan.Country,
		&beerScan.Price,
		&beerScan.Currency,
		&beerScan.CreatedAt,
		&beerScan.UpdatedAt,
	); errScan != nil {
		// Si no hay filas, devolvemos un error personalizado.
		if errors.Is(errScan, sql.ErrNoRows) {
			subLogger.Error().Msgf("no rows found: %v", errScan)
			return model.Beers{}, ErrBeerNotFound
		}
		subLogger.Error().Msgf("error scanning row: %v", errScan)
		return model.Beers{}, errScan
	}

	subLogger.Info().Msgf("FIN_OK | beer_id=%v", ID)
	return beerScan, nil
}

// CreateBeerWithId inserta una nueva cerveza en la base de datos, validando
// que el ID proporcionado no exista previamente.
//
// Parametros:
//   - ctx: contexto para control de tiempo de espera y cancelacion.
//   - beers: puntero a la instancia model.Beers que se va a insertar.
//
// Retorna:
//   - error: si el ID ya existe o hay fallo en la insercion/escaneo.
func (pb *beerRepository) CreateBeerWithId(ctx context.Context, beers *model.Beers) error {
	subLogger := log.With().Str("Method", "BeerRepository.CreateBeerWithId").Logger()
	subLogger.Info().Msg("INIT")
	subLogger.Info().Msgf("argument[s] beer_id=%v", beers.ID)

	// Preparar statement de insercion
	stmt, err := pb.Conn.PrepareContext(ctx, insertBeerWithId)
	if err != nil {
		subLogger.Error().Msgf("error preparing statement: %v", err)
		return err
	}
	// Asegurar cierre de statement
	defer func() {
		if errClose := stmt.Close(); errClose != nil {
			subLogger.Error().Msgf("error closing stmt: %v", errClose)
		}
	}()

	// Ejecutar QueryRowContext para obtener el ID generado o confirmar insercion
	row := stmt.QueryRowContext(ctx,
		&beers.ID,
		&beers.Name,
		&beers.Brewery,
		&beers.Country,
		&beers.Price,
		&beers.Currency,
		&beers.CreatedAt,
		&beers.UpdatedAt,
	)

	// Escanear el ID devuelto
	if errScan := row.Scan(&beers.ID); errScan != nil {
		subLogger.Error().Msgf("error scanning inserted ID: %v", errScan)
		return errScan
	}

	subLogger.Info().Msg("FIN_OK")
	return nil
}