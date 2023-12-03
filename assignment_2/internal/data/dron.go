package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agatai06/golang/internal/validator"
	"github.com/lib/pq"
	"time"
)

type Dron struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Price     Price     `json:"price,omitempty"`
	Materials []string  `json:"materials,omitempty"`
	Version   int32     `json:"version"`
}

func (p Dron) MarshalJSON() ([]byte, error) {
	var price string

	if p.Price != 0 {
		price = fmt.Sprintf("%d tenge", p.Price)
	}

	type ProductAlias Dron

	aux := struct {
		ProductAlias
		Price string `json:"price,omitempty"`
	}{
		ProductAlias: ProductAlias(p),
		Price:        price,
	}
	return json.Marshal(aux)
}

func ValidateProduct(v *validator.Validator, p *Dron) {
	v.Check(p.Title != "", "title", "must be provided")
	v.Check(len(p.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(p.Year != 0, "year", "must be provided")
	v.Check(p.Year >= 1888, "year", "must be greater than 1888")
	v.Check(p.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(p.Price != 0, "price", "must be provided")
	v.Check(p.Price > 0, "price ", "must be a positive integer")

	v.Check(p.Materials != nil, "categories", "must be provided")
	v.Check(len(p.Materials) >= 1, "categories", "must contain at least 1 category")
	v.Check(len(p.Materials) <= 5, "categories", "must not contain more than 5 categories")

	v.Check(validator.Unique(p.Materials), "categories", "must not contain duplicate values")
}

type DronModel struct {
	DB *sql.DB
}

func (p DronModel) Insert(product *Dron) error {
	//return nil
	query := `
		INSERT INTO products (title, year, price, categories)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{product.Title, product.Year, product.Price, pq.Array(product.Materials)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.CreatedAt, &product.Version)
}

func (p DronModel) Get(id int64) (*Dron, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, title, year, price, categories, version
		FROM products
		WHERE id = $1`

	var product Dron

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.Title,
		&product.Year,
		&product.Price,
		pq.Array(&product.Materials),
		&product.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (p DronModel) Update(product *Dron) error {
	query := `
		UPDATE products
		SET title = $1, year = $2, price = $3, categories = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []interface{}{
		product.Title,
		product.Year,
		product.Price,
		pq.Array(product.Materials),
		product.ID,
		product.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m DronModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM products
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

func (p DronModel) GetAll(title string, categories []string, filters Filters) ([]*Dron, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(),id, created_at, title, year, price, categories, version
		FROM products
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (categories @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, pq.Array(categories), filters.limit(), filters.offset()}
	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	products := []*Dron{}

	for rows.Next() {
		var product Dron
		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.CreatedAt,
			&product.Title,
			&product.Year,
			&product.Price,
			pq.Array(&product.Materials),
			&product.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		products = append(products, &product)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil

}
