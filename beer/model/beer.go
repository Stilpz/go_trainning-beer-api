// Package model define las estructuras de dominio y sus representaciones para la entidad Beer.
// Contiene tanto la definición interna de la cerveza en base de datos como su formato de respuesta HTTP.
package model

import "time"

// Beers representa la entidad Beer tal como se almacena en la base de datos.
// Cada campo está etiquetado para mapearse con la columna correspondiente.

type Beers struct {
	ID	uint `db:"id"`	//Identificador única de la cerveza
	Name	string `db:"name"`	//Nombre de la cerveza
	Brewery	string `db:"brewery"`	// Cervecería Productora
	Country	string `db:"country"`	// País de origen
	Price	float64 `db:"price"`	// Precio unitario de la cerveza en la moneda original
	Currency	string `db:"currency"`	// Código ISO de la moneda (ej: USD)
	CreatedAt	time.Time `db:"created_at"`	// Marca de tiempo de creación de la cerveza
	UpdatedAt	time.Time `db:"updated_at"`	// Marca de tiempo de última actualización de la cerveza
}

// BeersResponse define la estructura que se envía al cliente en las respuestas HTTP.
// Omite campos vacíos en el JSON (omitempty).
type BeersResponse struct {
	ID       uint    `json:"id,omitempty"`       // Identificador único de la cerveza
	Name     string  `json:"name,omitempty"`     // Nombre de la cerveza
	Brewery  string  `json:"brewery,omitempty"`  // Cervecería productora
	Country  string  `json:"country,omitempty"`  // País de origen
	Price    float64 `json:"price,omitempty"`    // Precio unitario
	Currency string  `json:"currency,omitempty"` // Código de moneda
	CreatedAt time.Time `json:"created_at,omitempty"` // Fecha de creación
	UpdatedAt time.Time `json:"updated_at,omitempty"` // Fecha de actualización
}

// PriceResponse representa la respuesta con el precio total de la caja de cervezas.
// PriceTotal contiene el monto calculado y se serializa como "price_total".
// CurrencyPay indica la moneda utilizada para el pago y se serializa como "currency_pay".
// Ambos campos se omiten del JSON si están vacíos o cero.
type PriceResponse struct {
	PriceTotal  float64 `json:"price_total,omitempty"`
	CurrencyPay string  `json:"currency_pay,omitempty"`
}

// ToBeersResponse transforma el modelo interno Beers a su representación para la respuesta HTTP.
// Devuelve un objeto BeersResponse con los campos públicamente expuestos.
func (b *Beers) ToBeersResponse() BeersResponse {
	// return BeersResponse{
	// 	ID:       b.ID,
	// 	Name:     b.Name,
	// 	Brewery:  b.Brewery,
	// 	Country:  b.Country,
	// 	Price:    b.Price,
	// 	Currency: b.Currency,
	// }
	    return BeersResponse{
        ID: b.ID,
        Name: b.Name,
        Brewery: b.Brewery,
        Country: b.Country,
        Price: b.Price,
        Currency: b.Currency,
        CreatedAt: b.CreatedAt,
        UpdatedAt: b.UpdatedAt,
    }
}

// BeersRequest define la estructura esperada en las peticiones HTTP
// para crear o actualizar una cerveza. Omite campos vacíos en el JSON.
type BeersRequest struct {
	ID       uint    `json:"id,omitempty"`       // Identificador único de la cerveza (opcional en creación)
	Name     string  `json:"name,omitempty"`     // Nombre de la cerveza
	Brewery  string  `json:"brewery,omitempty"`  // Cervecería productora
	Country  string  `json:"country,omitempty"`  // País de origen
	Price    float64 `json:"price,omitempty"`    // Precio unitario
	Currency string  `json:"currency,omitempty"` // Código ISO de la moneda (ej. "USD")
}

// Validate ejecuta validaciones de campo sobre la estructura BeersRequest.
// Retorna un mapa de errores donde la clave es el identificador de validación
// y el valor es el mensaje descriptivo.
// Errores posibles:
//   - "id_required": debe proporcionarse un ID válido (no cero).
//   - "name_required": el nombre no puede estar vacío.
//   - "brewery_required": la cervecería no puede estar vacía.
//   - "country_required": el país no puede estar vacío.
//   - "price_required": el precio debe ser distinto de cero.
//   - "currency_required": la moneda debe ser un código ISO válido de al menos 3 caracteres.
func (br *BeersRequest) Validate() map[string]string {
	errorMessages := make(map[string]string)

	if br.ID == 0 {
		errorMessages["id_required"] = "Id is required or id invalid"
	}
	if br.Name == "" {
		errorMessages["name_required"] = "name is required"
	}
	if br.Brewery == "" {
		errorMessages["brewery_required"] = "brewery is required"
	}
	if br.Country == "" {
		errorMessages["country_required"] = "country is required"
	}
	if br.Price == 0 {
		errorMessages["price_required"] = "price is required and must be greater than zero"
	}
	if br.Currency == "" || len(br.Currency) < 3 {
		errorMessages["currency_required"] = "currency is required and it has to be a valid currency code"
	}

	return errorMessages
}