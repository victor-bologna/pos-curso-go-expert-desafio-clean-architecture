package database

import (
	"database/sql"

	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/entity"
)

type OrderRepositoryInterface interface {
	Save(order *entity.Order) error
	GetAll() ([]entity.Order, error)
}

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (or *OrderRepository) Save(order *entity.Order) error {
	// Prepare stmt é usado para gerar várias execuções com o mesmo comando.
	stmt, err := or.DB.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetAll() ([]entity.Order, error) {
	rows, err := or.DB.Query("SELECT * FROM orders;")
	if err != nil {
		return nil, err
	}

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err = rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
