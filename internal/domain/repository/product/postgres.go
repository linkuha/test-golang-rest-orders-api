package product

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	productTableName = "products"
	pricesTableName  = "product_prices"
)

type repo struct {
	db *sql.DB
}

func newProductPostgresRepository(d *sql.DB) Repository {
	return &repo{
		db: d,
	}
}

func (r *repo) Get(ctx context.Context, id string) (*entity.Product, error) {
	query := fmt.Sprintf("SELECT id, name, description, left_in_stock FROM %s WHERE id = $1", productTableName)
	log.Debug().Msg("Query: " + query)

	row := r.db.QueryRowContext(ctx, query, id)
	product := entity.Product{}

	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.LeftInStock)
	if err != nil {
		return nil, errs.HandleErrorDB(err)
	}
	return &product, nil
}

func (r *repo) GetAll(ctx context.Context) (*[]entity.Product, error) {
	query := fmt.Sprintf("SELECT id, name, description, left_in_stock FROM %s", productTableName)
	log.Debug().Msg("Query: " + query)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errs.HandleErrorDB(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	products := []entity.Product{}
	for rows.Next() {
		p := entity.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.LeftInStock)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return &products, nil
}

func (r *repo) GetPrices(ctx context.Context, id string) (*[]entity.Price, error) {
	query := fmt.Sprintf("SELECT price, currency FROM %s WHERE product_id = $1", pricesTableName)
	log.Debug().Msg("Query: " + query)

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errs.HandleErrorDB(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	prices := []entity.Price{}
	for rows.Next() {
		p := entity.Price{}
		err := rows.Scan(&p.Price, &p.Currency)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		prices = append(prices, p)
	}

	return &prices, nil
}

func (r *repo) Store(ctx context.Context, product *entity.Product) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (name, description, left_in_stock) VALUES ($1, $2, $3) RETURNING id", productTableName)
	log.Debug().Msg("Query: " + query)

	row := r.db.QueryRowContext(ctx, query, product.Name, product.Description, product.LeftInStock)
	if err := row.Scan(&id); err != nil {
		return "", errs.HandleErrorDB(err)
	}

	return id, nil
}

func (r *repo) StoreWithPrices(ctx context.Context, product *entity.Product) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Debug().Msg("Start transaction err: " + err.Error())
		return "", errs.HandleErrorDB(err)
	}
	defer tx.Rollback()

	var productID string
	query := fmt.Sprintf("INSERT INTO %s (name, description, left_in_stock) VALUES ($1, $2, $3) RETURNING id", productTableName)
	log.Debug().Msg("Query: " + query)

	row := tx.QueryRowContext(ctx, query, product.Name, product.Description, product.LeftInStock)
	if err = row.Scan(&productID); err != nil {
		return "", errs.HandleErrorDB(err)
	}

	query = fmt.Sprintf(`INSERT INTO %s (product_id, currency, price) VALUES ($1, $2, $3)
		ON CONFLICT (product_id, currency) DO UPDATE SET price = EXCLUDED.price`, pricesTableName)
	log.Debug().Msg("Query for stmt: " + query)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Debug().Msg("Prepare stmt err: " + err.Error())
		return "", errs.HandleErrorDB(err)
	}

	for _, price := range product.Prices {
		_, err = stmt.ExecContext(ctx, productID, price.Currency, price.Price)
		if err != nil {
			return "", errs.HandleErrorDB(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Debug().Msg("Commit transaction err: " + err.Error())
		return "", errs.HandleErrorDB(err)
	}

	return productID, nil
}

func (r *repo) Update(ctx context.Context, id string, input *entity.ProductUpdateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.LeftInStock != nil {
		setValues = append(setValues, fmt.Sprintf("left_in_stock=$%d", argId))
		args = append(args, *input.LeftInStock)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", productTableName, setQuery, argId)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errs.HandleErrorDB(err)
	}

	return nil
}

func (r *repo) Remove(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", productTableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errs.HandleErrorDB(err)
	}
	return nil
}

func (r *repo) AddPrice(ctx context.Context, productID string, price *entity.Price) error {
	query := fmt.Sprintf(`INSERT INTO %s (product_id, currency, price) VALUES ($1, $2, $3)
		ON CONFLICT (product_id, currency) DO UPDATE SET price = EXCLUDED.price`, pricesTableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, productID, price.Currency, price.Price)
	if err != nil {
		return errs.HandleErrorDB(err)
	}
	return nil
}
