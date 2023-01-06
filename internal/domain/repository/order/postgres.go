package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	ordersTableName        = "user_orders"
	orderProductsTableName = "user_order_products"
	productsTableName      = "products"
)

type repo struct {
	db *sql.DB
}

func newOrderPostgresRepository(d *sql.DB) Repository {
	return &repo{
		db: d,
	}
}

func (r *repo) Get(ctx context.Context, id string) (*entity.Order, error) {
	query := fmt.Sprintf("SELECT id, user_id, number FROM %s WHERE id = $1", ordersTableName)
	log.Debug().Msg("Query: " + query)

	row := r.db.QueryRowContext(ctx, query, id)
	order := entity.Order{}

	if err := row.Scan(&order.ID, &order.UserID, &order.Number); err != nil {
		return nil, errs.HandleErrorDB(err)
	}
	return &order, nil
}

func (r *repo) GetAllByUserID(ctx context.Context, userID string) (*[]entity.Order, error) {
	query := fmt.Sprintf("SELECT id, user_id, number FROM %s WHERE user_id = $1", ordersTableName)
	log.Debug().Msg("Query: " + query)

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errs.HandleErrorDB(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	orders := []entity.Order{}
	for rows.Next() {
		o := entity.Order{}
		err := rows.Scan(&o.ID, &o.UserID, &o.Number)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		orders = append(orders, o)
	}

	return &orders, nil
}

func (r *repo) GetProducts(ctx context.Context, id string) (*[]entity.OrderProductView, error) {
	query := fmt.Sprintf("SELECT product_id, amount FROM %s WHERE order_id = $1", orderProductsTableName)
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

	products := []entity.OrderProductView{}
	for rows.Next() {
		p := entity.OrderProductView{}
		err := rows.Scan(&p.ID, &p.Amount)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return &products, nil
}

func (r *repo) Store(ctx context.Context, order *entity.Order) (string, error) {
	var id string
	// idempotent
	query := fmt.Sprintf(`WITH ins_orders AS (
    INSERT INTO %s (user_id, number)
    VALUES ($1, $2)
    ON CONFLICT(user_id, number) DO NOTHING
    RETURNING id
) SELECT COALESCE(
    (SELECT id FROM ins_orders),
    (SELECT id FROM %s WHERE user_id = $1 AND number = $2)
) as id`, ordersTableName, ordersTableName)
	log.Debug().Msg("Query: " + query)

	row := r.db.QueryRowContext(ctx, query, order.UserID, order.Number)
	if err := row.Scan(&id); err != nil {
		return "", errs.HandleErrorDB(err)
	}

	return id, nil
}

func (r *repo) Update(ctx context.Context, order *entity.Order) error {
	query := fmt.Sprintf("UPDATE %s SET number = $1, user_id = $2 WHERE id = $3", ordersTableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, order.Number, order.UserID, order.ID)
	if err != nil {
		return errs.HandleErrorDB(err)
	}

	return nil
}

func (r *repo) Remove(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ordersTableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errs.HandleErrorDB(err)
	}
	return nil
}

func (r *repo) AddProduct(ctx context.Context, op *entity.OrderProduct) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Debug().Msg("Start transaction err: " + err.Error())
		return errs.HandleErrorDB(err)
	}
	defer tx.Rollback()

	selQuery := fmt.Sprintf(`SELECT id, left_in_stock FROM %s WHERE id = $1`, productsTableName)
	log.Debug().Msg("Query: " + selQuery)

	row := tx.QueryRowContext(ctx, selQuery, op.ProductID)
	product := entity.Product{}

	if err = row.Scan(&product.ID, &product.LeftInStock); err != nil {
		return errs.HandleErrorDB(err)
	}

	// idempotent
	insQuery := fmt.Sprintf(`INSERT INTO %s (order_id, product_id, amount) VALUES ($1, $2, $3)
		ON CONFLICT (order_id, product_id) DO UPDATE SET amount = %s.amount + EXCLUDED.amount`, orderProductsTableName, orderProductsTableName)
	log.Debug().Msg("Query: " + insQuery)

	if _, err = tx.ExecContext(ctx, insQuery, op.OrderID, op.ProductID, op.Amount); err != nil {
		return errs.HandleErrorDB(err)
	}

	updateQuery := fmt.Sprintf(`UPDATE %s SET left_in_stock = $1 WHERE id = $2`, productsTableName)
	log.Debug().Msg("Query: " + updateQuery)

	if _, err = tx.ExecContext(ctx, updateQuery, product.LeftInStock-op.Amount, product.ID); err != nil {
		return errs.HandleErrorDB(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Debug().Msg("Commit transaction err: " + err.Error())
		return errs.HandleErrorDB(err)
	}

	return nil
}

func (r *repo) RemoveProduct(ctx context.Context, orderID, productID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE order_id = $1 AND product_id = $2`, orderProductsTableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, orderID, productID)
	if err != nil {
		return errs.HandleErrorDB(err)
	}

	return nil
}
