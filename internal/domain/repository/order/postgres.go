package order

import (
	"database/sql"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

const (
	ordersTableName        = "user_orders"
	orderProductsTableName = "user_order_products"
)

type repo struct {
	db *sql.DB
}

func newOrderPostgresRepository(d *sql.DB) Repository {
	return &repo{
		db: d,
	}
}

func (r *repo) Get(id string) (*entity.Order, error) {
	query := fmt.Sprintf("SELECT id, user_id, number FROM %s WHERE id = ?", ordersTableName)
	row := r.db.QueryRow(query, id)
	order := entity.Order{}

	err := row.Scan(&order.ID, &order.UserID, &order.Number)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *repo) GetAllByUserID(userID string) (*[]entity.Order, error) {
	query := fmt.Sprintf("SELECT id, user_id, number FROM %s WHERE user_id = ?", ordersTableName)
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
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
			fmt.Println(err)
			continue
		}
		orders = append(orders, o)
	}

	return &orders, nil
}

func (r *repo) GetProducts(id string) (*[]entity.OrderProductView, error) {
	query := fmt.Sprintf("SELECT product_id, amount FROM %s WHERE order_id = ?", orderProductsTableName)
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	var products []entity.OrderProductView
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

func (r *repo) Store(order *entity.Order) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (number, user_id) VALUES (?, ?) RETURNING id", ordersTableName)
	row := r.db.QueryRow(query, order.Number, order.UserID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Update(order *entity.Order) error {
	query := fmt.Sprintf("UPDATE %s SET number = ?, user_id = ? WHERE id = ?", ordersTableName)
	_, err := r.db.Exec(query, order.Number, order.UserID, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Remove(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", ordersTableName)
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddProduct(p *entity.OrderProduct) error {
	query := fmt.Sprintf(`INSERT INTO %s (order_id, product_id, amount) VALUES (?, ?, ?)
		ON CONFLICT (order_id, product_id) DO UPDATE SET amount = %s.amount + EXCLUDED.amount`, orderProductsTableName, orderProductsTableName)
	_, err := r.db.Exec(query, p.OrderID, p.ProductID, p.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) RemoveProduct(orderID, productID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE order_id = ? AND product_id = ?`, orderProductsTableName)
	_, err := r.db.Exec(query, orderID, productID)
	if err != nil {
		return err
	}

	return nil
}
