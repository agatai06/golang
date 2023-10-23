package data

import (
	"assignment_2/internal/validator"
	"encoding/json"
	"fmt"
	"time"
)

type Product struct {
	ID         int64     `json:"id"`                   // Unique integer ID for the product
	CreatedAt  time.Time `json:"-"`                    // Timestamp for when the product is added to our database
	Title      string    `json:"title"`                // Product title
	Year       int       `json:"year,omitempty"`       // Product release year
	Price      int       `json:"price,omitempty"`      // Product price (in tenge)
	Categories []string  `json:"categories,omitempty"` // Product category (technic, for home, etc.)
	Version    int32     `json:"version"`              // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

func (p Product) MarshalJSON() ([]byte, error) {
	// Create a variable holding the custom runtime string, just like before.
	var price string

	if p.Price != 0 {
		price = fmt.Sprintf("%d tenge", p.Price)
	}
	// Define a MovieAlias type which has the underlying type Movie. Due to the way that
	// Go handles type definitions (https://golang.org/ref/spec#Type_definitions) the
	// MovieAlias type will contain all the fields that our Movie struct has but,
	// importantly, none of the methods.
	type ProductAlias Product
	// Embed the MovieAlias type inside the anonymous struct, along with a Runtime field
	// that has the type string and the necessary struct tags. It's important that we
	// embed the MovieAlias type here, rather than the Movie type directly, to avoid
	// inheriting the MarshalJSON() method of the Movie type (which would result in an
	// infinite loop during encoding).
	aux := struct {
		ProductAlias
		Price string `json:"price,omitempty"`
	}{
		ProductAlias: ProductAlias(p),
		Price:        price,
	}
	return json.Marshal(aux)
}

func ValidateProduct(v *validator.Validator, p *Product) {
	v.Check(p.Title != "", "title", "must be provided")
	v.Check(len(p.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(p.Year != 0, "year", "must be provided")
	v.Check(p.Year >= 1888, "year", "must be greater than 1888")
	v.Check(p.Year <= int(time.Now().Year()), "year", "must not be in the future")

	v.Check(p.Price != 0, "price", "must be provided")
	v.Check(p.Price > 0, "price ", "must be a positive integer")

	v.Check(p.Categories != nil, "categories", "must be provided")
	v.Check(len(p.Categories) >= 1, "categories", "must contain at least 1 genre")
	v.Check(len(p.Categories) <= 5, "categories", "must not contain more than 5 genres")

	v.Check(validator.Unique(p.Categories), "categories", "must not contain duplicate values")
}
