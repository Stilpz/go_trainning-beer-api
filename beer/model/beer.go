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